//
//  AppDelegate.m
//  MazeRunner
//
//  Created by Delisa on 1/30/18.
//  Copyright © 2018 Bugsnag. All rights reserved.
//

#import "AppDelegate.h"
#import <Bugsnag/Bugsnag.h>

@interface InvariantException : NSException
@end

@implementation InvariantException
@end

@implementation AppDelegate


- (BOOL)application:(UIApplication *)application didFinishLaunchingWithOptions:(NSDictionary *)launchOptions {
    NSArray *launchArguments = [[NSProcessInfo processInfo] arguments];
    NSTimeInterval delay = 0;
    NSString *eventType = nil;
    NSString *mockAPIPath = nil;
    NSString *bugsnagAPIKey = nil;
    BOOL skipEvents = NO;
    for (NSString *argument in launchArguments) {
        if ([argument containsString:@"EVENT_DELAY"]) {
            delay = [[[argument componentsSeparatedByString:@"="] lastObject] integerValue];
        }
        if ([argument containsString:@"EVENT_TYPE"]) {
            eventType = [[argument componentsSeparatedByString:@"="] lastObject];
        }
        if ([argument containsString:@"MOCK_API_PATH"]) {
            mockAPIPath = [[argument componentsSeparatedByString:@"="] lastObject];
        }
        if ([argument containsString:@"BUGSNAG_API_KEY"]) {
            bugsnagAPIKey = [[argument componentsSeparatedByString:@"="] lastObject];
        }
        if ([argument containsString:@"SKIP_EVENTS"]) {
            skipEvents = YES;
        }
    }
    NSAssert(mockAPIPath != nil, @"The mock API path must be set prior to triggering events");

    BugsnagConfiguration *config = [BugsnagConfiguration new];
    config.apiKey = bugsnagAPIKey;
    config.notifyURL = [NSURL URLWithString:mockAPIPath];
    [Bugsnag startBugsnagWithConfiguration:config];

    NSLog(@"Arguments: %@", launchArguments);
    if (!skipEvents)
        [self triggerEventWithName:eventType delay:delay];

    return YES;
}

- (void)triggerEventWithName:(NSString *)name delay:(NSTimeInterval)delay {
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(delay * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
        if ([name isEqualToString:@"NSException"]) {
            @throw [InvariantException exceptionWithName:@"Invariant violation" reason:@"The cake was rotten" userInfo:nil];
        }
    });
}


@end
