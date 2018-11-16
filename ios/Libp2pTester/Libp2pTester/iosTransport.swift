//
//  iosTransport.swift
//  Libp2pTester
//
//  Created by why on 11/14/18.
//  Copyright © 2018 why. All rights reserved.
//

import UIKit
import Libp2p
import CocoaAsyncSocket

class iosTransport: NSObject, Libp2pTransportProtocol, GCDAsyncSocketDelegate {
    func protocols() -> Libp2pIntList! {
        let il = Libp2pIntList()
        il?.push(0x06)
        return il
    }
    
    func canDial(_ addr: Libp2pMultiaddr!) -> Bool {
        return true
    }
    
    func dial(_ raddr: Libp2pMultiaddr!, p1: Libp2pPeerID!) throws -> Libp2pConn {
        let mSocket = GCDAsyncSocket(delegate: self, delegateQueue: DispatchQueue.main)
        do {
            print("dialing...")
            
            var port:Int = 0
            try raddr.getPort(&port)
            
            try mSocket.connect(toHost: raddr.getHost(), onPort: UInt16(port))
            print("Socket connected!")
        } catch let error {
            print("socket connect: ", error)
        }
        throw MyError.BadThing("well, it might have worked")
    }
    
    enum MyError: Error {
        case BadThing(String)
    }
    func listen(_ laddr: Libp2pMultiaddr!) throws -> Libp2pListener {
        throw MyError.BadThing("not implemented")
    }
    
    func proxy() -> Bool {
        return false
    }

}
