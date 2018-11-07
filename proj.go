//SI413 Project 2 Bodeman, Sipes

package main

import "fmt"
import "math/rand"//https://gobyexample.com/random-numbers
import "strconv"
import "time"
import "os"
import "sync"
import "bufio"
//import "github.com/fatih/color" test this out for the printing out stuff


var randNum int;
var upTo int

type Player struct{ //player struct to save someone's data
  Name string //character name
  location int //character location
  notepad := [9]bool{false,false,false,false,false,false,false,false,false} //pieces of evidence
}


func saveCheck(){ //looks for save file
  //FIXME write this function
}

func setLocation(){ //Only sets the location of current user
  //TODO write this
}

func setUser(){ //read saved file and set User variables
  //TODO write this
}

func getLocation(){ //returns user location int
  //TODO write this
}

func prompt(loc int){ //print out approprate prompt for location and read in input
  //TODO write this
  fmt.Printf("%s\n", prompts[loc]); //prompts is a string array of location prompts
}

func execute(){
  var command string;
  var args string;
  while(true){
    fmt.Printf("> ");
    fmt.Scanf("%s",&command) //TODO this this reads in as expected
    reader := bufio.NewReader(os.Stdin);//https://stackoverflow.com/questions/14094190/golang-function-similar-to-getchar
    args, _ := reader.ReadString('\n');//reads the newline

    command = strings.ToLower(command);
    switch command {
    case "move":
      break;
    case "examine":
      break;
    case "ask" :
      break;
    case "accuse":
      break;
    case "look":
      break;
    case "notepad":
      break;
    case "save":
      break;
    case "quit":
      break;
    case "help":
    default:
      fmt.Println("Unknown command");//FIXME loop back to get command again
    }
  }

}

func game() {
  //TODO Make new user
  if(saveCheck()){
    fmt.Printf("Would you like to load the saved file? (Y/N)");
    var answer string;
    while(true){
      fmt.Scanf("%s", &answer);
      answer := strings.ToLower(answer[0]);
      if(strings.Compare(answer,"y") == 0){
        fmt.Printf("Loaded saved file\n");
        setUser()
        break;
      }else if(strings.Compare(answer, "n"){
        fmt.Printf("Starting new game\n");
        break;
      }else {
        fmt.Printf("Please specifiy Y or N\n");
      }
    }
  }
  while(true){
    prompt(getLocation());
    execute();
  }
}

func main() {
  go game()
  //Adds the current time as a seed for the random function
  //Provides synchronization functions for go routines (threads).

  var wg sync.WaitGroup //https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html
  //Indicates that there will be 2 threads that must be waited on. Sets the WaitGroup counter to 2.
  wg.Add(2)

  //Waits until the WaitGroup counter is 0.
  wg.Wait()
}
