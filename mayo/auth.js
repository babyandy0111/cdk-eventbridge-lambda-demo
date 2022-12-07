const msal = require('@azure/msal-node');

/**
 * Configuration object to be passed to MSAL instance on creation.
 * For a full list of MSAL Node configuration parameters, visit:
 * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-node/docs/configuration.md
 */
const msalConfig = {
	auth: {
		clientId: process.env.client_id,
		authority: process.env.aad_endpoint + process.env.tenant_id,
		clientSecret: process.env.client_secret,
	}
};

/**
 * With client credentials flows permissions need to be granted in the portal by a tenant administrator.
 * The scope is always in the format '<resource-appId-uri>/.default'. For more, visit:
 * https://docs.microsoft.com/azure/active-directory/develop/v2-oauth2-client-creds-grant-flow
 */
const tokenRequest = {
	scopes: [process.env.graph_endpoint + '.default'], // e.g. 'https://graph.microsoft.com/.default'
};

const apiConfig = {
	uri: process.env.graph_endpoint + 'v1.0/users', // e.g. 'https://graph.microsoft.com/v1.0/users'
	hr: process.env.graph_endpoint + 'v1.0/users/a811118a-0070-4640-a3aa-3c059685c9ba', // e.g. 'https://graph.microsoft.com/v1.0/users'
	andy: process.env.graph_endpoint + 'v1.0/users/70f2ae6b-9454-4941-a05a-adc216e38526', // e.g. 'https://graph.microsoft.com/v1.0/users'
};

/**
 * Initialize a confidential client application. For more info, visit:
 * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-node/docs/initialize-confidential-client-application.md
 */
const cca = new msal.ConfidentialClientApplication(msalConfig);

/**
 * Acquires token with client credentials.
 * @param {object} tokenRequest
 */
async function getToken(tokenRequest) {
	return await cca.acquireTokenByClientCredential(tokenRequest);
}

module.exports = {
	apiConfig: apiConfig,
	tokenRequest: tokenRequest,
	getToken: getToken
};