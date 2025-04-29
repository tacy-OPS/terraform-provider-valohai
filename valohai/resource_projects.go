package valohai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_notifications": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {

	// URL de l'API Valohai pour créer un projet
	url := "https://app.valohai.com/api/v0/projects/"

	// Récupérer la configuration du provider
	config := m.(map[string]interface{})
	authToken := config["token"].(string)

	// Validation des champs obligatoires
	if d.Get("name").(string) == "" {
		return fmt.Errorf("Le champ 'name' est obligatoire")
	}
	if d.Get("owner").(string) == "" {
		return fmt.Errorf("Le champ 'owner' est obligatoire")
	}

	// Préparer le payload
	payload := map[string]interface{}{
		"name":  d.Get("name").(string),
		"owner": d.Get("owner").(string),
	}

	// Champs optionnels
	if v, ok := d.GetOk("description"); ok {
		payload["description"] = v.(string)
	}
	if v, ok := d.GetOk("template_url"); ok {
		payload["template"] = v.(string)
	}
	if v, ok := d.GetOk("default_notifications"); ok {
		payload["default_notifications"] = v.(string) == "true"
	}

	// Encodage JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'encodage du payload: %w", err)
	}

	// Création de la requête HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Erreur lors de la création de la requête: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+authToken)

	// Logs pour le débogage
	fmt.Printf("Requête POST envoyée à l'URL : %s\n", url)
	fmt.Printf("Payload : %s\n", string(body))

	// Envoi de la requête
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout de 10 secondes
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'envoi de la requête: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("Erreur lors de la fermeture du corps de la réponse: %s\n", cerr)
		}
	}()

	// Vérification du code HTTP
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Non autorisé (401) : vérifiez votre token d'authentification")
	} else if resp.StatusCode == http.StatusBadRequest {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("Requête invalide (400) : %v", errResp)
	} else if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Erreur API %d : %s", resp.StatusCode, resp.Status)
	}

	// Décodage de la réponse
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("Erreur lors du décodage de la réponse: %w", err)
	}

	// Enregistrer l'ID du projet dans Terraform
	d.SetId(result.ID)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	// Implement logic to read a widget.
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {

	// Récupérer la configuration du provider
	config := m.(map[string]interface{})
	authToken := config["token"].(string)

	// Récupérer l'ID du projet à modifier
	projectID := d.Id()

	// Construire l'URL de l'API Valohai pour modifier un projet
	url := fmt.Sprintf("https://app.valohai.com/api/v0/projects/%s/", projectID)

	// Préparer le payload
	payload := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}

	// Encodage JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'encodage du payload: %w", err)
	}

	// Création de la requête HTTP
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Erreur lors de la création de la requête: %w", err)
	}

	// Ajouter le token d'autorisation dans les headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+authToken)

	// Ajouter un timeout au client HTTP
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout de 10 secondes
	}

	// Effectuer la requête
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'envoi de la requête PUT: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("Erreur lors de la fermeture du corps de la réponse: %s\n", cerr)
		}
	}()

	// Logs pour le débogage
	fmt.Printf("Requête PUT envoyée à l'URL : %s\n", url)
	fmt.Printf("Code de statut de la réponse : %d\n", resp.StatusCode)

	// Vérifier le code de statut de la réponse
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Projet introuvable (404)")
	} else if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Non autorisé (401) : vérifiez votre token d'authentification")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Erreur lors de la modification du projet: statut %s", resp.Status)
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {

	// Récupérer la configuration du provider
	config := m.(map[string]interface{})
	authToken := config["token"].(string)

	// Récupérer l'ID du projet à supprimer
	projectID := d.Id()

	// Construire l'URL de l'API Valohai pour supprimer un projet
	url := fmt.Sprintf("https://app.valohai.com/api/v0/projects/%s/", projectID)

	// Effectuer la requête HTTP DELETE
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Erreur lors de la création de la requête DELETE: %s", err)
	}

	// Ajouter le token d'autorisation dans les headers
	req.Header.Set("Authorization", "Token "+authToken)

	// Ajouter un timeout au client HTTP
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout de 10 secondes
	}

	// Effectuer la requête
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'envoi de la requête DELETE: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("Erreur lors de la fermeture du corps de la réponse: %s\n", cerr)
		}
	}()

	// Logs pour le débogage
	fmt.Printf("Requête DELETE envoyée à l'URL : %s\n", url)
	fmt.Printf("Code de statut de la réponse : %d\n", resp.StatusCode)

	// Vérifier le code de statut de la réponse
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Projet introuvable (404)")
	} else if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Non autorisé (401) : vérifiez votre token d'authentification")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Erreur lors de la suppression du projet: statut %s", resp.Status)
	}

	// Si la suppression est réussie, supprimer l'ID de la ressource dans Terraform
	d.SetId("") // Cela marque la ressource comme supprimée dans Terraform
	return nil
}
