#!/usr/bin/env python3.6
# -*- coding: utf-8 -*-
"""Sample code for registering a file in an account in the PencilDATA server"""

import base64, ssl, json, hashlib, sys#Imports all the standard libraries neeeded
import getpass
import requests
import json
if sys.version_info[0] < 3:
    Input = raw_input
else:
    Input = input

sslcontext=ssl.create_default_context()#Create a SSL context

#We define the functions we will use:
def PencilDATA_login(username, password):
    """Construct an ProvenanceValidator object by logging in to the
    Pencildata server.
       Both the username and the password arguments may be given as str.
       Password bytes sequence will be submitted to the server encoded in
    base64. After a successful authentication, the login_data property is
    populated by a dictionary that (among other things) contains the
    user's IdToken. If the authentication fails, an exception is raised
    (actually a HTTPError: Bad Request)."""

    url = 'https://api.pencildata.com/token'
    data = {'userId': username, 'password': password}
    head = {"Content-Type": "application/json"}
    res = requests.request("POST", url, data=json.dumps(data), headers=head)

    return res.json()

def register_file(login_data,file,storage="none"):
    """Register a file (by its SHA-256 hash) in your Pencildata account.
    Warning: this method does not check if the file hash exists in the
    registers. It returns the asset id for the file.

    Arguments:
    file: file name or a file object. If file is given as a file-like
    object, this method advances the current position of the file until
    its end, but it does not close the file-like object
    storage: 'public' or 'private'. Whether to store the file entry in the
    public or in the private database at the PencilDATA server."""

    datajson = {}
    datajson["hash"] = file_hash(file)
    datajson["storage"] = storage
    url = "https://api.pencildata.com/register/"

    head = {"Content-Type": "application/json","Authorization": "Bearer {0}".format(login_data['data']['accessToken'])}#Request HTTP headers
    res = requests.request("POST", url, data=json.dumps(datajson), headers=head)

    return res.json()
def file_hash(file):
    """Returns the hash for the file, as it is sent to the Pencildata
    server for registration or for verification
    Arguments:
        file: The access path to the file to calculate the hash for."""

    hashobj=hashlib.sha256()

    with open(file, "rb") as fileobj:
        chunk=fileobj.read(65536)
        while len(chunk)>0:
            hashobj.update(chunk)
            chunk=fileobj.read(65536)

    return hashobj.hexdigest()

#Then, we log into the PencilDATA server, using the account and
#its password. If the login is successful, returns a
#pencildata.ProvenanceValidator object, which holds the user identification tokens.
#If login fails (either because the password is incorrect, the user name does not
#exist or there is a connection problem), it raises an exception.
username = Input("username: ")
password = getpass.getpass("password: ")
pencillogindata=PencilDATA_login(username,password)

#The following statement registers the file SHA256 hash in the PencilDATA server, in the
#account whose data is stored in the pencildata obj, and return the asset id that
#the server returned. Note than the file is never uploaded to the server; the only
#data sent to the PencilDATA server is the hexadecimal digest of the SHA256 hash
#for the file, which is calculated locally. For this statement, we submit as
#arguments the login data dictionary we obtained in the last step, and the
#name of the file to register, in this order.
asset_id=register_file(pencillogindata,"sample.txt")
print("The asset id of the newly registered file is: ",asset_id)


