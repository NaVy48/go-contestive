# Pusiau automatinis būdas sugeneruoti klasių diagramą
echo '@startuml' > 'cd.uml';
# Klasių generavimas
# Skriptas praleidžia anonimines klases
# 511-oje eilutėje reikia pakeisti rankiniu būdu
pcregrep -M 'type +[A-Z|a-z]+ +struct +{({}|[^}])*}' `find .. -name *.go` \
  | sed -e 's:{}:():' \
        -e 's:\t: :g' \
        -e 's:/\([a-z|A-Z|_|0-9]\+\)\(\.go\)\?:namespace \1\2 {\n:g' \
        -e 's:\.\.::' -e 's:\: *type \([A-Z|a-z|_|0-9]*\) struct:class \1:' \
  | awk ' BEGIN { c=0 }
          /namespace/ { c++ } 
          /}/ { for (i = 0; i < c; i++) { print "}"  }  c=0 } 
          {print $0}' \
  | sed -e 's://.*$::' \
        -e 's:`.*$::' \
        -e 's: \+: :g' \
        -e 's:\.go:_go:' \
  | grep -v '^ *$' \
  >> 'cd.uml';
# Sąsajų generavimas
# Jokių klaidų nesugeneruoja
pcregrep -M 'type +[A-Z|a-z]+ +interface +{({}|[^}])*}' `find .. -name *.go` \
  | sed -e 's:{}:():' \
        -e 's:\t: :g' \
        -e 's:/\([a-z|A-Z|_|0-9]\+\)\(\.go\)\?:namespace \1\2 {\n:g' \
        -e 's:\.\.::' \
        -e 's:\: *type \([A-Z|a-z|_|0-9]*\) interface:interface \1:' \
  | awk ' BEGIN { c=0 }
          /namespace/ { c++ }
          /}/ { for (i = 0; i < c; i++) { print "}"  }  c=0 }
          {print $0}' \
  | sed -e 's://.*$::' \
        -e 's:`.*$::' \
        -e 's: \+: :g' \
        -e 's:\.go:_go:' \
  | grep -v '^ *$' \
  >> 'cd.uml';
echo '@enduml' >> 'cd.uml';
echo 'Pakeisti anoniminę struktūrą 511 eilutėje rankiniu būdu ir paleisti:

  env PLANTUML_LIMIT_SIZE=24576 plantuml cd.uml

Gautos klasių diagramos peržiūra:
  
  exo-open cd.png
';

