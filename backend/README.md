
# Backend-dokumentation

## Nätverk
Backend använder sig av TCP sockets, port 1337

## Kommunikation
För att kommunicera med servern, uppräta en anslutning med TCP på port 1337.

När anslutningen är upprättad används följande gränssnitt (interface)

Servern kommer koppla samman en anslutning med användaren när den loggar in, vilket gör att användaren inte behöver skicka med sin användare i alla anrop

### Format
Servern vill ha en string med formatet "\*\*anropsfunktion, arg1, arg2, .., argn\*\*"


### joinroom 
Argument: roomName

Ansluter användaren till ett rum

Exempel: "joinroom, room1"

### enterroom
Går in i rummet


### exitroom
Argument: roomname

Tar bort användaren från ett rum
