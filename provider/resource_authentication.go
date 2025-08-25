package provider

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/terraform-provider-iis/iis"
)

func resourceAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthenticationCreate,
		ReadContext:   resourceAuthenticationRead,
		UpdateContext: resourceAuthenticationUpdate,
		DeleteContext: resourceAuthenticationDelete,

		Schema: map[string]*schema.Schema{
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"anonymous": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"user": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"basic": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"default_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"realm": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"windows": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
						"providers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceAuthenticationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	application := d.Get("application").(string)
	tflog.Debug(ctx, "Creating authentication: "+toJSON(application))
	auth, err := client.ReadAuthenticationFromApplication(ctx, application)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := updateAuthProviders(ctx, d, client, auth); err != nil {
		return err
	}
	tflog.Debug(ctx, "Created authentication: "+toJSON(auth))
	d.SetId(auth.ID)
	return nil
}

func resourceAuthenticationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	auth, err := client.ReadAuthentication(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read authentication: "+toJSON(auth))
	if err = readAuthenticationProvider(ctx, d, "anonymous", buildAnonymousAuthProvider(client, &auth)); err != nil {
		return diag.FromErr(err)
	}
	if err = readAuthenticationProvider(ctx, d, "basic", buildBasicAuthProvider(client, &auth)); err != nil {
		return diag.FromErr(err)
	}
	if err = readAuthenticationProvider(ctx, d, "windows", buildWindowsAuthProvider(client, &auth)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(auth.ID)
	return nil
}

func resourceAuthenticationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	tflog.Debug(ctx, "Updating authentication: "+toJSON(d.Id()))
	auth, err := client.ReadAuthentication(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return updateAuthProviders(ctx, d, client, auth)
}

func resourceAuthenticationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func updateAuthProviders(ctx context.Context, d *schema.ResourceData, client *iis.Client, auth iis.Authentication) diag.Diagnostics {
	anonymousAuthProvider := buildAnonymousAuthProvider(client, &auth)
	basicAuthProvider := buildBasicAuthProvider(client, &auth)
	windowsAuthProvider := buildWindowsAuthProvider(client, &auth)

	if err := updateAuthenticationProvider(ctx, d, client, "anonymous", anonymousAuthProvider, updateAnonymousAuthentication); err != nil {
		return diag.FromErr(err)
	}
	if err := updateAuthenticationProvider(ctx, d, client, "basic", basicAuthProvider, updateBasicAuthentication); err != nil {
		return diag.FromErr(err)
	}
	if err := updateAuthenticationProvider(ctx, d, client, "windows", windowsAuthProvider, updateWindowsAuthentication); err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Updated authentication: "+toJSON(auth))

	return nil
}

func updateAnonymousAuthentication(ctx context.Context, client *iis.Client, auth interface{}, data map[string]interface{}) error {
	anonymous := auth.(iis.AnonymousAuthentication)
	anonymous.Enabled = data["enabled"].(bool)
	anonymous.User = data["user"].(string)

	_, err := client.UpdateAnonymousAuthentication(ctx, &anonymous)

	return err
}

func updateBasicAuthentication(ctx context.Context, client *iis.Client, auth interface{}, data map[string]interface{}) error {
	basic := auth.(iis.BasicAuthentication)
	basic.Enabled = data["enabled"].(bool)
	basic.DefaultLogonDomain = data["default_domain"].(string)
	basic.Realm = data["realm"].(string)

	_, err := client.UpdateBasicAuthentication(ctx, &basic)

	return err
}

func updateWindowsAuthentication(ctx context.Context, client *iis.Client, auth interface{}, data map[string]interface{}) error {
	enabledProviders := data["providers"].([]interface{})
	windows := auth.(iis.WindowsAuthentication)
	windows.Enabled = data["enabled"].(bool)
	for i, provider := range windows.Providers {
		windows.Providers[i].Enabled = false
		for _, name := range enabledProviders {
			if strings.EqualFold(provider.Name, name.(string)) {
				windows.Providers[i].Enabled = true
				break
			}
		}
	}

	_, err := client.UpdateWindowsAuthentication(ctx, &windows)

	return err
}

func updateAuthenticationProvider(ctx context.Context, d *schema.ResourceData, client *iis.Client, key string, fetch FetchAuthProvider, update UpdateAuthProvider) error {
	if !d.HasChange(key) {
		log.Printf("no changes for %s authentication", key)
		return nil
	}
	provider, err := fetch(ctx)
	if err != nil {
		return err
	}

	if !hasNestedMap(d, key) {
		return nil
	}

	return update(ctx, client, provider, getNestedMap(d, key))
}

func buildAnonymousAuthProvider(client *iis.Client, auth *iis.Authentication) FetchAuthProvider {
	return func(ctx context.Context) (AuthProvider, error) {
		return client.ReadAnonymousAuthentication(ctx, auth)
	}
}

func buildBasicAuthProvider(client *iis.Client, auth *iis.Authentication) FetchAuthProvider {
	return func(ctx context.Context) (AuthProvider, error) {
		return client.ReadBasicAuthentication(ctx, auth)
	}
}

func buildWindowsAuthProvider(client *iis.Client, auth *iis.Authentication) FetchAuthProvider {
	return func(ctx context.Context) (AuthProvider, error) {
		return client.ReadWindowsAuthentication(ctx, auth)
	}
}

type AuthProvider interface {
	ToMap() map[string]interface{}
}
type FetchAuthProvider func(ctx context.Context) (AuthProvider, error)

type UpdateAuthProvider func(context.Context, *iis.Client, interface{}, map[string]interface{}) error

func readAuthenticationProvider(ctx context.Context, d *schema.ResourceData, key string, fetch FetchAuthProvider) error {
	provider, err := fetch(ctx)
	if err != nil {
		return err
	}
	providerList := make([]map[string]interface{}, 1)
	providerList = append(providerList, provider.ToMap())
	if err := d.Set(key, providerList); err != nil {
		return err
	}
	return nil
}
