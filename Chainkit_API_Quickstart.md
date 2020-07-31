
# Chainkit API Quickstart

Last updated: November 14, 2019

## Summary:
This quickstart covers the four easy steps to using the Chainkit API: Hash, Authenticate, Register, Verify


## Step 1: Generate a Hash
Create a digital fingerprint, a *hash value*, of your file or data. We have 
[sample code](https://github.com/chainkit/QuickStart/tree/master/sample-code)
for many programming languages and platforms.

### Hash a File

Mac/Linux Example:

```bash
cat sample-document.pdf | openssl sha256
```

Javascript Example:

```javascript
crypto.createHmac('sha256', fileContents)
```

Java Example:

```java
hashString(fileContents, "SHA-256");
```

### Output Value

The *Hash Value* generated for that specific file would be the same on every platform.
For example:

```
ba3f474830169ddaece741cbbfbec13086139d88907104564b765c3935109d64
```

## Step 2: Authenticate to Chainkit API

Submit your username and password to the Chainkit token API. [Get an account here](https://chainkit.com/start) if needed.
The API returns JSON with multiple values. You will need to extract the accessToken value for use in the Next Step.

**Example Request:**

```bash
curl -s -X POST https://api.chainkit.com/token -H 'Content-Type: application/json' -d '{"userId":"your-username", "password": "your-password"}'
```

**Example Response:**

```json
{"data":{"expiresIn":"3600", "accessToken": "eyJraWQiOiJETzRxQWFGb...", "refreshToken": "..."}}
```

In the response above, you need to copy the full accessToken `eyJraWQiOiJETzRxQWFGb...` which we have truncated in order to save space.

## Step 3: Register Your Hash Value

Submit the *hash value* you created in Step 1 along with the `accessToken` which you created in Step 2 to the Chainkit Register API. 

If successful, the API will returns an **assetId** unique to that registration and hash value. Save this **assetId**, as it will
be used in Step 4 to verify the file.

Chainkit allows storage on private and public blockchains. For example, you could choose to store your hash on the
public Ethereum blockchain, or a private blockchain including Hyperledger Sawtooth or VMware Blockchain.
Private is faster and recommended for testing.

In this example we are storing the *hash value* on Chainkit's Hyperledger Sawtooth blockchain. Please contact us
at <support@chainkit.com> if you would want to upgrade your account permissions to use other blockchains. 

**Example Request:**

Here is the basic syntax of the **Register() API**.

```bash
curl -s -X POST https://api.chainkit.com/register -H 'Content-Type: application/json' -d '{"hash":"your-hash-value", "storage":"pencil"}' -H 'Authorization: Bearer 'your-accessToken-value''
```

So, following along in this example, substitute your values from the prior steps for `"your-hash-value"`
and `"your-accessToken-value"`. Note: Pay close attention to the single quote marks that must be including
 at the end of the command.

```bash
curl -s -X POST https://api.chainkit.com/register -H 'Content-Type: application/json' -d  '{"hash":"ba3f474830169ddaece741cbbfbec13086139d88907104564b765c3935109d64", "storage":"pencil"}' -H 'Authorization: Bearer 'eyJraWQiOiJETzRxQWFGb...''
```

**Example Response** 

```json
{"assetId":"1526398776829"}
```

Save the **assetId** for later verification.

## Step 4: Verify Your Hash Value

To verify that your file has not been tampered with, first create a **Fresh Hash**, using the same hashing formula from Step 1. **Do not reuse**
the original hash you used at registration, as that defeats the purpose.

Then you will submit the fresh hash, the **assetId** you got from registration, and your **accessToken**
to the **Verify API**. It will return *true* if your new hash matches what was originally submitted to us
or *false* to indicate the hash value does not match and something in your file or data has changed.

**Example Request:**

Here is the basic syntax of the Verify() API.

```bash
curl -s -X GET 'https://api.chainkit.com/verify/your-assetId?hash=your-freshly-generated-hash-value&storage=private' -H 'Authorization: Bearer your-accessToken-value'
```

So, in this example, you will need to subsitute your values from prior steps for `"assetID"` and 
`"your-freshly-generated-hash-value"` and `"your-accessToken-value"`. Note: Pay close attention
to the single quote marks that must be including at the end of the command.

```bash
curl -s -X GET 'https://api.chainkit.com/verify/1526398776829?hash=ba3f474830169ddaece741cbbfbec13086139d88907104564b765c3935109d64&storage=pencil' -H 'Authorization: Bearer 'eyJraWQiOiJETzRxQWFGb...''
```

**Example Response**

```json
{"verified":true}
```
This response shows that your fresh hash matches the original stored on the Chainkit API for that **assetId**.

Note: You can test this verification by altering the original file, hashing the altered version, and
submitting that to the verify API with the original **assetId**. It will return false, since the hash
of the altered file does not match what was originally registered with us.

# Need help?
Contact us at [info@chainkit.com](mailto:info@chainkit.com)
