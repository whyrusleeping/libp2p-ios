//
//  ViewController.swift
//  Libp2pTester
//
//  Created by why on 11/13/18.
//  Copyright © 2018 why. All rights reserved.
//

import UIKit
import Libp2p

class ViewController: UIViewController {
    
    var host:Libp2pHost? = nil

    @IBOutlet weak var theLabel: UILabel!
    @IBOutlet weak var peerIDbox: UITextField!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        let casTpt = iosTransport()
        
        var opterr:NSError?
        host = Libp2pNew(casTpt, &opterr)
        if let err = opterr {
            print("bad error: ", err)
            return
        }
        
        peerIDbox.text = host!.peerInfo().id_().string()
        
        let pinfo = Libp2pParseMultiaddrString("/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", &opterr)
        if let err = opterr {
            print("bad error: ", err)
            return
        }
        
        do {
            try host!.connect(pinfo)
        } catch {
            print("connect failed: \(error)")
        }
        
    }
    
    @IBAction func clickTheButton(_ sender: Any) {
        let start = DispatchTime.now()
        doThePing()
        let end = DispatchTime.now()
        
        let nanoTime = end.uptimeNanoseconds - start.uptimeNanoseconds // <<<<< Difference in nano seconds (UInt64)
        let timeInterval = Double(nanoTime) / 1_000_000_000
        theLabel.text = "\(timeInterval) seconds"
    }
    
    func doThePing() {
        print("doing the ping")
        
        // parse this again because globals suck
        var opterr:NSError?
        let pinfo = Libp2pParseMultiaddrString("/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", &opterr)
        if let err = opterr {
            print("bad error: ", err)
            return
        }
        
        do {
            let stream = try host!.newStream(pinfo?.id_(), proto: "/ipfs/ping/1.0.0")
            var data = Data(count: 32)
            data[4] = 6 // just so i can distinguish it from an empty array
            
            var n:Int = 0
            try stream.write(data, ret0_: &n)
            
            print("Wrote data: \(n)")
            
            let recv = try stream.readData(32)
            
            let str = String(data: recv.base64EncodedData(), encoding: String.Encoding.utf8)
            print("I think we pinged! \(str)")

            
            try stream.close()
        } catch {
            print("new stream failed: \(error)")
        }
    }
    
}

