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
The proposed endpoints follow.

Upload a DICOM file
-------------------
POST /upload ::
    
    Content-Type: application/json
    {
        "file": path of the file to upload
    }

Response ::

    200 OK
    Content-Type: application/json
    {
        "fileID" : "1234"
    }

    400 Bad Request
    Content-Type: application/json
    {
        "Message": "File missing or not valid"
    }


Retrieve a header attribute by tag
----------------------------------
GET /fileID/attributes?tag=tag ::

    200 OK
    Content-Type: application/json
    {
        "Tag": tag,
        "Value Representation": VR
    }

GET /fileID/attributes?tag=tag ::

    404 Not Found

GET /fileID/attributes?tag=tag ::

    400 Bad Request
    Content-Type: application/json
    {
        "Message": "Incorrect or malformed tag"
    }


Convert the DICOM file into a PNG file
--------------------------------------
GET /fileID/png ::

    200 OK
    Content-Type: image/png
    Binary representation of the PNG file

GET /fileID/png ::
    
    404 Not Found
