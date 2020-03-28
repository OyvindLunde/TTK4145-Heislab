## TODO's fremover:

- Finne en måte for heisene å sette egen id
- Når og hvordan elevlist settes og oppdateres
- Timer: Watchdog, timeout osv
- Kun ta en ordre når en annen heis har bekreftet den (kun bekrefte ting som kommer fra nettverket?)
- Fikse kommunikasjon
- Globale variabler: Hvordan lage, dele, endre osv.

- Prioritize CAB orders? Might fuck up costFunction thou. 
- Mulig løsning: Når "ferdig i execute, se først etter cab orders før man går til IDLE". Dette for å unngå "trolls" som får heisen f.eks. til å kjøre opp og ned mellom 1. og 2. etg hele tiden. Løser ikke nødvendigvis problemer for feks. ned 4 etg, men vet heller ikke helt hva oppgaven krever.
- Fikse State er laget to steder


### Litt lenger frem i tid
- UpdateLights() må endres slik at den ikke setter lampa før ordren er confirmed. 
- Er også en unødvendig overlapp mellom UpdateLights() og HandleButtonEvents(), som begge setter lys.