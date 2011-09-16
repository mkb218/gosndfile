#include "virtual.h"

sf_count_t  gocall_get_filelen (void *user_data) {
	return gsfLen(user_data);
}

sf_count_t  gocall_seek (sf_count_t offset, int whence, void *user_data) {
	return gsfSeek(offset, whence, user_data);
}

sf_count_t  gocall_read        (void *ptr, sf_count_t count, void *user_data) {
	return gsfRead(ptr, count, user_data);
}

sf_count_t  gocall_write       (const void *ptr, sf_count_t count, void *user_data) {
	return gsfWrite(ptr, count, user_data);
}

sf_count_t  gocall_tell        (void *user_data) {
	return gsfTell(user_data);
}

SF_VIRTUAL_IO* virtualio() {
	SF_VIRTUAL_IO *svi = malloc(sizeof(SF_VIRTUAL_IO));
	svi->get_filelen = gocall_get_filelen;
	svi->seek = gocall_seek;
	svi->read = gocall_read;
	svi->write = gocall_write;
	svi->tell = gocall_tell;
	return svi;
}
	