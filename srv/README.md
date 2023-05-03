# Srv

Srv is a lightweight Go server designed for serving static files.
It is intended to be used to serve a web UI (e.g. React app) *behind your load balancer*.

It assumes:
* TLS has already been terminated
* Authentication/authorisation has already been performed

Srv supports `http/1.1` and `h2c` (plain-text `http/2`) which it will transparently use based on incoming requests.
It also uses the [AutoKubeOps Serverless Framework](https://gitlab.com/autokubeops/serverless) and can therefore be run as a Serverless function.

## Configuration

`SRV_STATIC_DIR`: path to the static files

`SRV_DOT_ENV`: path to the `.env` file

`SRV_ENV_FILE`: name of the generated `.js` file to be placed within `SRV_STATIC_DIR`

### Dotenv support

Srv has built in capability to generate runtime variables for Web applications.
Any variables defined in a dotenv file that are found in the Srv environment will be imported at startup.

1. Create a `.env` file containing key=value data for the environment variables you wish to use.
   If you do not wish to provide a value, key= can be used.
    ```dotenv
    API_URL=https://example.org
    API_KEY=
    ```
2. Enable the feature by setting `SRV_DOT_ENV` to the name of the file created above
3. Consume the variables from your application
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