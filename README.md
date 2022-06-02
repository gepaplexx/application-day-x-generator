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

- Grundgerüst eines ValueBuilders aus Vorlagen/XXXValueBuilder.go kopieren. Eine neue Datei in src/pkg/generator anlegen. Bsp.: 
```go
// Datei pgk/src/generator/MyNewValueBuilder.go
package generator

import (
	"gepaplexx/day-x-generator/pkg/util"
)

type MyNewValueBuilder struct{}

func (gen *MyNewValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	// TODO

	return values, nil
}
```

- GetValues implementieren. Bsp.:
```go
// Datei pgk/src/generator/MyNewValueBuilder.go
package generator

import (
	utils "gepaplexx/day-x-generator/pkg/util"
)

type MyNewValueBuilder struct{}

func (gen *MyNewValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	values["key1"] = config["val1"]
    values["key2"] = utils.Base64(config["key2"])

	return values, nil
}
```
- Neuen Generator aktivieren: Dazu muss ein neuer Eintrag in src/pkg/generator/Configuration.go gemacht werden. Bsp.:
```go
// Datei src/pkg/generator/Configuration.go
...
var GENERATORS = []Generator{
	{
		ValueBuilder: &MyNewValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "my-new-value-builder",
	},
}
...
```

## Verwendung

- Vorlagen/config.yaml kopieren und für die jeweilige Umgebung ausfüllen
- Ausführen des Generators (Annahme: augefülltes config.yaml liegt im aktuellen Verzeichnis):
    - Via go: 
        - Nach "src" wechseln
        - ```go run . config.yaml```
    - Via docker: 
        - ```docker run --rm -it -v $(pwd)/config.yaml:/app/config.yaml -v $(pwd)/generated:/app/generated ghcr.io/gepaplexx/day-x-generator config.yaml```
- Im Verzeichnis "generated" befinden sich die generierten Value-Files und Key/Zertifikat zum entschlüsseln der SealedSecrets
- Hinweis: Existiert bereits ein Verzeichnis “generated” und enthält dieses einen Private Key und eine Zertifikat für die jeweilige Umgebung ([env].crt und [env.key]) werden diese verwendet und keine neuen generiert.