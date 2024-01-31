package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"todo_task": resourceTask(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Retrieve the API URL from the provider configuration.
	url := d.Get("url").(string)

	// Normally, you would establish a connection to your backend service here.
	return &ProviderConfig{URL: url}, nil
}

type ProviderConfig struct {
	URL string
}

func resourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskCreate,
		ReadContext:   resourceTaskRead,
		UpdateContext: resourceTaskUpdate,
		DeleteContext: resourceTaskDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_done": &schema.Schema{
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

	// Call your API or perform the necessary logic to create the task.
	// taskID := YourAPICreateFunction(config.URL, name, isDone)

	taskID := 123 // Replace this with the actual task ID obtained from your API.

	d.SetId(strconv.Itoa(taskID))
	return resourceTaskRead(ctx, d, m)
}

func resourceTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	taskID, _ := strconv.Atoi(d.Id())

	// Call your API or perform the necessary logic to read the task details.
	// taskDetails := YourAPIReadFunction(config.URL, taskID)

	taskDetails := map[string]interface{}{
		"name":    "Task Name", // Replace this with the actual task name.
		"is_done": false,       // Replace this with the actual task status.
	}

	d.Set("name", taskDetails["name"])
	d.Set("is_done", taskDetails["is_done"])

	return nil
}

func resourceTaskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	name := d.Get("name").(string)
	isDone := d.Get("is_done").(bool)

	// Call your API or perform the necessary logic to update the task.
	// YourAPIUpdateFunction(config.URL, d.Id(), name, isDone)

	return resourceTaskRead(ctx, d, m)
}

func resourceTaskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)

	// Call your API or perform the necessary logic to delete the task.
	// YourAPIDeleteFunction(config.URL, d.Id())

	return nil
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}

func resourceTaskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	name := d.Get("name").(string)
	isDone := d.Get("is_done").(bool)

	// Call your API or perform the necessary logic to create the task.
	apiURL := fmt.Sprintf("%s/tasks", config.URL) // Assuming your API endpoint for creating tasks is "/tasks".

	requestBody := fmt.Sprintf(`{"name": "%s", "is_done": %t}`, name, isDone)

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(requestBody))
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to create task. Status code: %d", resp.StatusCode)
	}

	// Parse the response to get the created task ID.
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

	// Call your API or perform the necessary logic to read the task details.
	apiURL := fmt.Sprintf("%s/tasks/%d", config.URL, taskID) // Assuming your API endpoint for reading a task is "/tasks/{taskID}".

	resp, err := http.Get(apiURL)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to read task. Status code: %d", resp.StatusCode)
	}

	// Parse the response to get the task details.
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

	// Call your API or perform the necessary logic to update the task.
	apiURL := fmt.Sprintf("%s/tasks/%d", config.URL, taskID) // Assuming your API endpoint for updating a task is "/tasks/{taskID}".

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

	// Call your API or perform the necessary logic to delete the task.
	apiURL := fmt.Sprintf("%s/tasks/%d", config.URL, taskID) // Assuming your API endpoint for deleting a task is "/tasks/{taskID}".

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

