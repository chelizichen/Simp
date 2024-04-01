import { useShellStore } from "@/stores/counter";

export function NewEventSource(command:string ){
    const BASE = '/simpexpansionserver/shell/input?command='+command;
    const eventSource = new EventSource(BASE);
    const store = useShellStore()
    eventSource.onmessage = function(event) {
      const eventData = JSON.parse(event.data);
      console.log('Received event:', eventData);
      if(eventData.message == 'done'){
        return eventSource.close()
      }
      store.pushStack(
        `<span style="color:white;font-size:16px;font-weight:600;">${eventData.message}</span>`
        )
    };
    eventSource.onerror = function(error) {  
      if (eventSource.readyState !== EventSource.CLOSED) {  
          eventSource.close();  
      }  
      console.log('error',error);
      
  };  
}
