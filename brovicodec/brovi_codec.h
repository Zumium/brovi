#ifndef _BROVI_CODEC
#define _BROVI_CODEC
#define BROVI_CODEC_ENCODE_ERR -1
#define BROVI_CODEC_ZERO_SIZE_ERR -2

#ifdef __cplusplus
extern "C" {
#endif

typedef struct
{
    int width;
    int height;
} BroviCodecConfig;

typedef struct
{
    void *data;
    int size;
} H264Frame;

typedef void BroviCodec;

BroviCodec *BroviCodec_New(BroviCodecConfig);
void BroviCodec_Close(BroviCodec *);
int BroviCodec_EncodeFrame(BroviCodec *, void *, H264Frame *);
int BroviCodec_FlushDelayedFrame(BroviCodec *, H264Frame *);

#ifdef __cplusplus
}
#endif

#endif