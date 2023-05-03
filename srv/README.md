# Srv

Srv is a lightweight Go server designed for serving static files.
It is intended to be used to serve a web UI (e.g. React app) *behind your load balancer*.

It assumes:
* TLS has already been terminated
* Authentication/authorisation has already been performed

Srv supports `http/1.1` and `h2c` (plain-text `http/2`) which it will transparently use based on incoming requests.
It also uses the [AutoKubeOps Serverless Framework](https://gitlab.com/autokubeops/serverless) and can therefore be run as a Serverless function.

## Configuration

* `NIB_DATA_PATH`: path to the static files
* `NIB_ENV_FILE`: name of the generated `.js` file to be placed within `NIB_DATA_PATH`. Defaults to `env-config.js`

### Dotenv support

Srv has built in capability to generate runtime variables for Web applications.
Any variables defined in a dotenv file that are found in the Srv environment will be imported at startup.

1. Create a `.env` file containing key=value data for the environment variables you wish to use.
   If you do not wish to provide a default value, key= can be used.
    ```dotenv
    API_URL=https://example.org
    API_KEY=
    ```
2. Consume the variables from your application
    ```html
      <!DOCTYPE html>
    <html lang="en">
    <head>
        <script src="env-config.js"></script>
           <script src="index.js"></script>
        <title>Hello, World!</title>
    </head>
    </html>
   ```
   ```javascript
   // usage from javascript
   console.log(window._env_.API_URL);
   ```
   ```typescript
   // typescript declarations
   declare global {
       interface Window {
           _env_?: { 
               API_URL?: string; 
               API_KEY?: string; 
           }
       }
   }
   ```