using UnityEngine;
using System.Collections;
using System.Collections.Generic;

public class StreamCData  {

	private Dictionary<string,object> header;
	private Dictionary<string,object> body;

	public StreamCData(){}

	public void SetHeader(Dictionary<string,object> d){
		this.header = d;
	}
	public void SetBody(Dictionary<string,object> d){
		this.body = d;
	}

	public string GetCMD(){
		return (string)this.header["CMD"];
	}

	public Dictionary<string,Dictionary<string,object>> GetData(){
		Dictionary<string,Dictionary<string,object>> data = new Dictionary<string,Dictionary<string,object>>();
		data["H"] = this.header;
		data["B"] = this.body;
		return data;
	}

}
