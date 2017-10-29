#ifndef _BROVI_CAM
#define _BROVI_CAM
#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct
{
    char *devfile;
    int width;
    int height;
} BroviCamConfig;

typedef struct
{
    void *start;
    size_t length;
} VideoBuffer;

typedef void BroviCam;

BroviCam *BroviCam_Open(BroviCamConfig *);
void BroviCam_Close(BroviCam *);

#ifdef __cplusplus
}
#endif

#endif
