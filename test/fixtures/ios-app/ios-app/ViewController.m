//
//  ViewController.m
//  MazeRunner
//
//  Created by Delisa on 1/30/18.
//  Copyright © 2018 Bugsnag. All rights reserved.
//

#import "ViewController.h"

@interface ViewController ()
@property (weak, nonatomic) IBOutlet UITextField *urlField;
@end

@implementation ViewController
- (IBAction)triggerNSException:(id)sender {
    NSURL *notifyURL = [NSURL URLWithString:self.urlField.text];
    NSLog(@"Notifying %@", notifyURL);
}

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view, typically from a nib.
}


- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}


@end
