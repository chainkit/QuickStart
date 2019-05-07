#!/usr/bin/env python3.6
# -*- coding: utf-8 -*-
"""Sample code for verifying the provenance of a a file in the PencilDATA server"""
import base64, urllib, ssl, json, hashlib#Imports all the standard libraries neeeded
import requests, sys, getpass
sslcontext=ssl.create_default_context()#Create a SSL context
def PencilDATA_login(username,password):
    """Construct an ProvenanceValidator object by logging in to the 
    Pencildata server.
       Both the username and the password arguments may be given as str. 
       Password bytes sequence will be submitted to the server encoded in 
    base64. After a successful authentication, the login_data property is 
    populated by a dictionary that (among other things) contains the 
    user's accessToken. If the authentication fails, an exception is raised 
    (actually a HTTPError: Bad Request)."""

    url = 'https://api.pencildata.com/token'
    datajson = {'userId': username, 'password': password}
    head = {"Content-Type": "application/json"}
    res = requests.request("POST", url, data=json.dumps(datajson), headers=head)
    return res.json()
def file_hash(file):
    """Returns the hash for the file, as it is sent to the Pencildata 
    server for registration or for verification
    Arguments:
        file: The access path to the file to calculate the hash for."""

    hashobj=hashlib.sha256()
    with open(file,"rb") as fileobj:
        chunk=fileobj.read(65536)
        while len(chunk)>0:
            hashobj.update(chunk)
            chunk=fileobj.read(65536)
    return hashobj.hexdigest()
def verify_file(login_data, asset_id, file, storage="none"):
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
    # datajson = "{0}".format(datajson).encode("utf-8")
    url = "https://api.pencildata.com/verify/"+asset_id
    head = {"Content-Type": "application/json","Authorization": "Bearer {0}".format(login_data['data']['accessToken'])} #Request HTTP headers
    res = requests.request("GET", url, params=datajson, headers=head)
    return res.json()

if __name__ == '__main__':
    #Then, we log into the PencilDATA server, using the account and
    #its password. If the login is successful, returns a 
    #pencildata.ProvenanceValidator object, which holds the user identification tokens.
    #If login fails (either because the password is incorrect, the user name does not 
    #exist or there is a connection problem), it raises an exception.
    filename = sys.argv[-1]
    username = input("username: ")
    password = getpass.getpass("password: ")
    asset_id = input("assetId: ")
    storage = input("storage: ")
    pencillogindata=PencilDATA_login(username, password)
    #In a run in the demo account, we registered the sample.pdf file under the demo account 
    #and got the asset id XXXXXXX. You should change the asset id variable to 
    #the asset id you want to use, or change the parameters directly in the veryfy_file 
    #method calls (this variable exists here only to clarify that all the following 
    #verify_file method calls use the same value as the first input)
    #Verifies the provenance of the same file we used to register in PencilDATA.
    #For doing this, we use the veryfy_file method, using as first argument the login 
    #data dictionary obtained in the last step, next the asset id that we got when 
    #we registered the hash of the file we want to test provenance for,
    #and last the name of the file (its absolute path or its path relative to
    #the current directory)
    versamefile=verify_file(pencillogindata, asset_id, filename, storage)#this should return True
    print("Verification for the file used to register the asset: ", versamefile)
    #Verifies the provenance of a modified copy of the file we used to register in PencilDATA
    vercopymod1=verify_file(pencillogindata, asset_id, "modified-sample.txt", storage)#This should return False
    print("Verification for a modified copy of the file used to register the asset: ", vercopymod1)
