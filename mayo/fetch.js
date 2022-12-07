const axios = require('axios');

async function callApi(endpoint, accessToken) {
    const options = {
        headers: {
            Authorization: `Bearer ${accessToken}`
        }
    };

    console.log('request made to web API at: ' + new Date().toString());

    try {
        const response = await axios.default.get(endpoint, options);
        return response.data;
    } catch (error) {
        console.log(error)
        return error;
    }
};

async function callEventApi(endpoint, sdate, edate, accessToken) {
    const options = {
        headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json'
        }
    }

    const data = {
        "subject": "Leave List",
        "start": {
            "dateTime": sdate,
            "timeZone": "UTC"
        },
        "end": {
            "dateTime": edate,
            "timeZone": "UTC"
        },
        "body": {
            "contentType": "HTML",
            "content": "context～～～～～～"
        },
        "attendees": [
            {
                "emailAddress": {
                    "address": "team@test.com",
                    "name": "TEST CORP"
                },
                "type": "required"
            }
        ],
    }

    try {
        const response = await axios.default.post(endpoint, data, options)
        return response.data;
    } catch (error) {
        // console.log(error)
        return error;
    }
};

module.exports = {
    callApi: callApi,
    callEventApi: callEventApi
};
