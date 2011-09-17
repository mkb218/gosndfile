#ifndef GOSNDFILE_VIRTUAL
#define GOSNDFILE_VIRTUAL

#include <sndfile.h>

sf_count_t  gocall_get_filelen (void *user_data) ;
sf_count_t  gocall_seek        (sf_count_t offset, int whence, void *user_data) ;
sf_count_t  gocall_read        (void *ptr, sf_count_t count, void *user_data) ;
sf_count_t  gocall_write       (const void *ptr, sf_count_t count, void *user_data) ;
sf_count_t  gocall_tell        (void *user_data) ;

SF_VIRTUAL_IO *virtualio();

#endif