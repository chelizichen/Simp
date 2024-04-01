export function NewEventSource(command:string ){
    const BASE = '/simpexpansionserver/shell/input?';
    const Search = new URLSearchParams()
    Search.set('command',command)
    const Target = BASE + Search.toString();
    console.log('v',Target);
    const eventSource = new EventSource(BASE);
    eventSource.onmessage = function(event) {
      const eventData = JSON.parse(event.data);
      console.log('Received event:', eventData);
    };
}
