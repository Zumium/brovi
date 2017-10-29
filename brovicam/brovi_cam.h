#ifndef _BROVI_CAM
#define _BROVI_CAM
#include <stddef.h>
#include <linux/videodev2.h>

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

typedef struct
{
    VideoBuffer *buffer;
    struct v4l2_buffer v4l_buf;
} VideoBufferStatus;

typedef void BroviCam;

BroviCam *BroviCam_Open(BroviCamConfig *);
void BroviCam_Close(BroviCam *);
int BroviCam_Start(BroviCam *);
int BroviCam_Stop(BroviCam *);

VideoBufferStatus BroviCam_NextBufferA(BroviCam *);
int BroviCam_NextBufferB(BroviCam*,VideoBufferStatus);

#ifdef __cplusplus
}
#endif

#endif
