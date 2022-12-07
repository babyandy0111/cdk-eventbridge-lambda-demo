#!/usr/outlook/env node

require('dotenv').config();
const fetch = require('./fetch');
const auth = require('./auth');
const {Client} = require("@microsoft/microsoft-graph-client");
require("isomorphic-fetch");

exports.handler = async function (event, context, callback) {
    console.log(process.env.mayo_api_domain)
    console.log(process.env.mayo_key)
    context.callbackWaitsForEmptyEventLoop = false;

    const dateObj = new Date();
    const month = dateObj.getUTCMonth() + 1; //months from 1-12
    const day = dateObj.getUTCDate();
    const year = dateObj.getUTCFullYear();

    const sdate = year + "/" + month + "/" + day + 'T09:00:00.000Z';
    const edate = year + "/" + month + "/" + day + 'T18:00:00.000Z';

    try {
        const authResponse = await auth.getToken(auth.tokenRequest);
        const subscription = auth.apiConfig.hr + "/calendar/events"
        const event = await fetch.callEventApi(subscription, sdate, edate, authResponse.accessToken);
        console.log(event);

    } catch (error) {
        console.log(error);
    }

    return {
        statusCode: 200,
        headers: {"Content-Type": "text/json"},
        body: JSON.stringify({message: "Hello from my Lambda node!"})
    };
};
