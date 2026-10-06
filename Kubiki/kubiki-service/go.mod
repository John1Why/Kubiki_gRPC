module kubiki_Service

go 1.25

require (
	github.com/John1Why/proto_nri_kubiki v0.0.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx/v5 v5.7.6
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

replace github.com/John1Why/proto_nri_kubiki => ../../Proto_kubiki/api