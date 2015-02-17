using UnityEngine;
using System.Collections;
using System.Net.Sockets;
using System.Threading;
using MsgPack;
using MsgPack.Serialization;
using System.Collections.Generic;
using System.IO;


public class StreamClient {

	private TcpClient client;
	private bool threadOK = true;

	public StreamClient(){}

	public void Connect(string host,int port){
		this.client = new TcpClient(host,port);

		Thread thread = new Thread( () =>Receiver());
		thread.Start ();
		

	}
	
	public void Close(){
		this.threadOK  = false;
		if(this.client != null ){

			NetworkStream stream = this.client.GetStream();
			if (stream != null) {
				stream.Close ();
				stream = null;
			}
			this.client.Close ();
			this.client = null;
		}
		Debug.Log ("Closed");
	}

	public void Send(StreamCData cdata){

		MemoryStream stream = new MemoryStream();
		var serializer = MessagePackSerializer.Get<Dictionary<string, Dictionary<string, object>>>();
		serializer.Pack(stream, cdata.GetData ());

		this.client.GetStream().Write(stream.GetBuffer(),0,(int)stream.Length);

	}

	private void Receiver(){

		while(this.threadOK ){

			while (this.client.Available == 0)
			{
				Thread.Sleep(10);
			}

			try
			{
				MessagePackObject mpo = Unpacking.UnpackObject( this.client.GetStream() );
				Debug.Log(mpo);

			}
			catch( UnpackException  )
			{
				Thread.Sleep(10);
			}

			
		}
	}

}


