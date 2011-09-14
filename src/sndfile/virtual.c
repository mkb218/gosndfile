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
