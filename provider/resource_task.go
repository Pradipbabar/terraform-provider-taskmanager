package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskCreate,
		ReadContext:   resourceTaskRead,
		UpdateContext: resourceTaskUpdate,
		DeleteContext: resourceTaskDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_done": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceTaskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	name := d.Get("name").(string)
	isDone := d.Get("is_done").(bool)

	apiURL := fmt.Sprintf("%s/tasks", config.URL)
	requestBody := fmt.Sprintf(`{"name": "%s", "is_done": %t}`, name, isDone)

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to create task. Status code: %d", resp.StatusCode)
	}

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return diag.FromErr(err)
	}

	taskID := int(responseData["id"].(float64))
	d.SetId(strconv.Itoa(taskID))

	return resourceTaskRead(ctx, d, m)
}

func resourceTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	taskID, _ := strconv.Atoi(d.Id())
	apiURL := fmt.Sprintf("%s/tasks/%d", config.URL, taskID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to read task. Status code: %d", resp.StatusCode)
	}

	var taskDetails map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&taskDetails); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", taskDetails["name"])
	d.Set("is_done", taskDetails["is_done"])

	return nil
}

func resourceTaskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	taskID, _ := strconv.Atoi(d.Id())
	name := d.Get("name").(string)
	isDone := d.Get("is_done").(bool)

	apiURL := fmt.Sprintf("%s/tasks/%d", config.URL, taskID)
	requestBody := fmt.Sprintf(`{"name": "%s", "is_done": %t}`, name, isDone)

	req, err := http.NewRequest("PUT", apiURL, strings.NewReader(requestBody))
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to update task. Status code: %d", resp.StatusCode)
	}

	return resourceTaskRead(ctx, d, m)
}

func resourceTaskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	taskID, _ := strconv.Atoi(d.Id())
	apiURL := fmt.Sprintf("%s/tasks/%d", config.URL, taskID)

	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to delete task. Status code: %d", resp.StatusCode)
	}

	return nil
}
