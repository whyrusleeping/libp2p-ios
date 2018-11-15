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
    
    var libp2p:Libp2pLibp2p? = nil

    @IBOutlet weak var theLabel: UILabel!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        
        var opterr:NSError?
        let libp2p = Libp2pNew(&opterr)
        if let err = opterr {
            print("bad error: ", err)
            return
        }
        
        let pinfo = Libp2pParseMultiaddrString("/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", &opterr)
        if let err = opterr {
            print("bad error: ", err)
            return
        }
        
        do {
            try libp2p?.connect(pinfo)
        } catch {
            print("connect failed: \(error)")
        }
        
        doThePing()
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
        
        // parse this again because globals suck
        var opterr:NSError?
        let pinfo = Libp2pParseMultiaddrString("/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", &opterr)
        if let err = opterr {
            print("bad error: ", err)
            return
        }
        
        do {
            let stream = try libp2p?.newStream(pinfo?.id_(), proto: "/ipfs/ping/1.0.0")
            var data = Data(count: 32)
            data[4] = 6
            
            let str1 = String(data: data.base64EncodedData(), encoding: String.Encoding.utf8)
            print("about to send message \(str1)")
            
            var n:Int = 0
            try stream?.write(data, ret0_: &n)
            
            print("Wrote data: \(n)")
            
            let recv = try stream?.readData(32)
            
            if let readdata = recv {
                let str = String(data: readdata.base64EncodedData(), encoding: String.Encoding.utf8)
                print("I think we pinged! \(str)")
            } else {
                print("possibly failed to ping?", recv)
            }
            
            try stream?.close()
        } catch {
            print("new stream failed: \(error)")
        }
    }
    
}

