SRC_DIR=../proto
EXT_SRC_DIR=../components/protobuf-ext/proto

proto:
	protoc -I=${SRC_DIR} --python_out=. ${SRC_DIR}/tracelet_location.proto ${SRC_DIR}/tracelet_status.proto
