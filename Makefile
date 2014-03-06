GO=go

SOURCE=m2u.go
EXEC=m2u

all:
	${GO} build ${SOURCE}

clean:
	rm -f ${EXEC}
