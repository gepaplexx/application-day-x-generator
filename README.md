# application-day-x-generator

Dieses Repository enthält ein Tool zum Generieren der values.yaml-Dateien, welche zum Aufsetzen eines neuen Clusters benötigt werden.

## Erweiterung

Um einen neuen Generator zu implementieren sind folgende Schritte notwendig:

- In Vorlagen/config.yaml im entsprechenden Schritt (initial-cluster-setup, cluster-setup-checkpoint, cluster-applications) einen neuen Eintrag anlegen, Bsp.: 

```yaml
# Datei Vorlagen/config.yaml
...
cluster-applications:
    my-new-value-builder:
        key1: val1
        key2: val2
...
```
- Template der Zieldatei erweitern (src/tempates/[initial-cluster-setup | cluster-setup-checkpoint | cluster-applications]). Bsp.:
```yaml
# Datei src/templates/cluster-applications.yaml.tpl
...
      helm:
        parameters:
          - name: "my.new.values.1"
            value: "{{ .key1 }}"
          - name: "my.new.values.2.base64"
            value: "{{ .key2 }}"
...
```
- Generator aktivieren: Dazu muss ein neuer Eintrag in src/pkg/generator/Configuration.go gemacht werden. In den meisten Fällen reicht der GenericCopyValueBuilder aus. Bsp.:
```go
// Datei src/pkg/generator/Configuration.go
...
var GENERATORS = []Generator{
    {
        ValueBuilder: &GenericCopyValueBuilder{},
        Stage:        ClusterApplications,
        Name:         "my-new-value-builder",
    },
}
...
```

## Verwendung

- Vorlagen/config.yaml kopieren und für die jeweilige Umgebung ausfüllen
- Ausführen des Generators:
    - Via go: 
        - Nach "src" wechseln
        - ```go run . /path/to/config.yaml```
    - Via docker (Annahme: augefülltes config.yaml liegt im aktuellen Verzeichnis): 
        - ```docker run --rm -it -v $(pwd)/config.yaml:/app/config.yaml -v $(pwd)/generated:/app/generated ghcr.io/gepaplexx/day-x-generator config.yaml```
- Im Verzeichnis "generated" befinden sich die generierten Value-Files und Key/Zertifikat zum entschlüsseln der SealedSecrets
- **Hinweis:** Existiert bereits ein Verzeichnis “generated” und enthält dieses einen Private Key und ein Zertifikat für die jeweilige Umgebung ([env].crt und [env.key]) werden diese verwendet und keine neuen generiert.