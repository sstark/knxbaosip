
knxbaosip
=========

Simple Go client implementation of the KNX IP BAOS application layer web
interface, as specified for accessing the following home automation devices
from Weinzierl:

  - KNX IP BAOS 771
  - KNX IP BAOS 772
  - KNX IP BAOS 773
  - KNX IP BAOS 774
  - KNX IP BAOS 777 

Original API documentation:

  - https://www.weinzierl.de/images/download/documents/baos/KNX_IP_BAOS_WebServices.pdf

Notes
=====

low level functions
-------------------

registerDatapoint(int, type)

get...json  (raw json)
get...raw   (raw bytes)
get...dpt   (carries dpt type in struct)


get...dpt:

  - get description first
  - then use json api call corresponding to dpt type
  - return struct with
    - original json
    - go api type corresponding to dpt type
    - methode to coerce type of variable

mid level functions
-------------------

getAs...int (coerced getter)

getAs...int
getAs...float
getAs...date
getAs...string

  - forces value into given type
  - warns if embedded dpt type does not match


high level functions
--------------------

showDescription()
set()
toggle()

