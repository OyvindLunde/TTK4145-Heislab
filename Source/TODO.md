## TODO's fremover:

- Finne en måte for heisene å sette egen id // Kanskje, ikke et krav
- Timer: Watchdog, timeout osv
- Kun ta en ordre når en annen heis har bekreftet den (kun bekrefte ting som kommer fra nettverket?)


- Fikse State er laget to steder
- Ordrer går av og på ^ - må fikses. Fikset et problem, men fortsatt ikke 100 %. Lys er ikke good, og det ser av og til ut som shit blocker hvis man spammer ordrer (full channel?)
- Splitt opp ShouldItakeOrder i to funksjoner, detect og solve

- Alt av feilhåndtering

- ordne slik at tickeren trigger et interupt

Spør studass om:
- Sleep time, særlig for network
- Polletallet på simulator
- main moduler og sånt
- Evt heiser på forskjellige pcer (ikke så viktig)



## Endringslogg
- Lite tillegg i updateOrderList som fikset orreproblemet (måtte bare sjekke at Finished==false, er fordi vi sletter ordren først etter 3 sek).
- Endret Display til å kun vise egen aktive ordre til å være grønn, dvs at kun en knapp er grønn per heis
- La til checkForUpdates() for å også oppdatere displayet når en annen heis (from network) har endringer. Økte poll rate på updates() i Display.go til 100msek for å passe på at Displayet ikke kræsjer. Kræsjer fortsatt vel full spam thou, må fikse Release() ?
- Endret displayet til å ha "passe" størrelse ved oppstart 



Kom tilbake til orderHandler: StopAtFloor()
Gjøre om status (og mer?) til enum/struct
Slett sleep i caborderbackup hvis good
Se over InitMyElevInfo()