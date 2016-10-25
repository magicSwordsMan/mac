#ifndef driver_h
#define driver_h

#import <Cocoa/Cocoa.h>

@interface DriverDelegate : NSObject <NSApplicationDelegate>
@property NSMenu *dock;

- (instancetype)init;
@end

void *Driver_Init();
void Driver_Run();
void Driver_Terminate();

#endif /* driver_h */