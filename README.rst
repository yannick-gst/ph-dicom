=======
phdicom
=======

**phdicom** is an internal microservice intended to accept and store an uploaded DICOM file, extract and return any DICOM header attribute based on a DICOM Tag as
a query parameter, and finally convert the file into a PNG for browser-based viewing. It is written in Go.

DICOM
=====
DICOM is a standard that is used for the digital storage and transmission of medical images. More information about the DICOM standard can be found `here <https://en.wikipedia.org/wiki/DICOM>`_.

API Endpoints
=============
The proposed endpoints follow, with examples of possible outputs.

Upload a DICOM file
-------------------
POST /upload ::
    
    Content-Type: application/json
    {
        "file": path of the file to upload
    }

Assumption: given it's an internal microservice, the server will know of the file path at the time of the request.

Response ::

    200 OK
    Content-Type: application/json
    {
        "fileID" : "1234"
    }

Assumption: it's okay to return the fileID as a string.

    400 Bad Request
    Message: "File missing or not valid"


Retrieve a header attribute by tag
----------------------------------
GET /file/fileID/attributes?tagGroup=group&tagElement=element ::

    200 OK
    Content-Type: application/json
    Response returned verbatim from dicom

GET /file/fileID/attributes?tagGroup=group&tagElement=element ::

    404 Not Found

GET /file/fileID/attributes?tagGroup=group&tagElement=element ::

    400 Bad Request
    Message: "Incorrect or malformed tag"


Convert the DICOM file into a PNG file
--------------------------------------
GET /file/fileID/png ::

    200 OK
    Content-Type: image/png
    Binary representation of the PNG file

GET /file/fileID/png ::
    
    404 Not Found

Building and running
====================
To build, run `go build`. To test, run `go test`

Future improvements
===================
In the interest of submitting a solution as quickly as possible, I did not implement the PNG conversion part.
This should be added in the future.

The retrieving by attribute endpoint does not currently define an actual representation of the attributes.

There are no failure test cases as of now. Those would be needed for more robust testing.