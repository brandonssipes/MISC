//SI413 Project - Part 2 
//Bodeman, Sipes

package main

import "fmt"
import "math/rand"
import "strings"
import "time"
import "os"
import "os/exec"
import "runtime"
import "sync"
import "bufio"

import "strconv"
import "os/signal"//signal handling https://stackoverflow.com/questions/18106749/golang-catch-signals
import "syscall"

import "crypto/sha256"

/*
 * Information needed to track the user's progress through the game  
 */
type Player struct {
  name string       //character name
  location int      //current location
  notepad [12]bool  //found evidence
  visited [4]bool   //key locations
}

//Initializes user as global variable with all values in struct set to the type's equivalent of 0;
var user Player;
//Determines if a given professor is currently in their office
var present int;
//channel to tell threads to quit safely
//var quitting int = 0 //= make(chan bool, 1)
quitting := make(chan bool);

func checkSave() bool{ //looks for save file
  if _, err := os.Stat("./.Murder"); os.IsNotExist(err){//https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
    return false;
  }else{
    return true;
  }
  return false;
}


func setUser(){ //read saved file and set User variables
  fd, _ := os.Open("./.Murder");
  defer fd.Close();

  scanner := bufio.NewScanner(fd);
  scanner.Scan()
  user.name = scanner.Text()

  scanner.Scan()
  data := scanner.Text()
  user.location,_ = strconv.Atoi(data)//https://golang.org/pkg/strconv/

  scanner.Scan()
  notes := strings.Split(scanner.Text(), ",")

  for i:=0; i < 12; i++{
    user.notepad[i],_ = strconv.ParseBool(notes[i])
  }
  scanner.Scan()
  visit := strings.Split(scanner.Text(), ",")
  for i:=0; i < 4; i++{
    user.visited[i], _ = strconv.ParseBool(visit[i])
  }
}


/*
 * Checks the current location of the user and prints out the appropriate scenario
 */
func prompt(){

  loc := user.location;

  //https://stackoverflow.com/questions/7933460/how-do-you-write-multiline-strings-in-go
  //Golang doesn't like random line breaks in a string, so I concatenated the strings just for source code readability

  //Intro Scene
  if loc == 0 {
    fmt.Printf("You wake up with a start as a deafening pounding shakes the door to your room. Your heart is racing as you lie frozen in bed," +
    " confused and disoriented. You hear the steady patter of rain against the window as a voice drifts through a crack in your door.\n\n\"%s," +
    " we have a situation. A body was just found, impaled on the USNA crest infront of Michelson. They are requesting your presence at the" +
    " murder scene as the lead investigator. Get dressed quickly. I'll be waiting for you there.\"\n\nYou are uncertain who the voice belongs" +
    " too, but after being shocked out of your slumber, you slip out of bed and hastily pull on your uniform.\n\nType \"Help\" for a list of" +
    " actions. To continue type \"Move Crime Scene\".\n\n", user.name);
  }

  //Intro Crime Scene
  if loc == 1 {
    fmt.Printf("The morning air has a slight chill to it and you pull your jacket closer to your body to preserve some warmth. As you come upon" +
    " Michelson, you notice the police slowly removing a body from where it was impaled on the USNA crest which is attached above the main entrance" +
    " to the building. Blood trickles down the front steps and the first rays of light cast the scene in an eerie glow. A man standing at the crime" +
    " scene notices you walking over and quickly makes his way towards you.\n\n\"%s, I'm Detective Holmes with Homicide. The Superintendant notified" +
    " us that you are the Naval Academy's best hope at finding who the killer is. Our resources lead us to believe that the killer is someone who" +
    " works inside Michelson. Why don't you take a \"Look\" around. You are free to investigate anywhere inside the building and when you are ready" +
    " to accuse someone, come find me.\n\n", user.name);

    //This scenario should only print. Ensures that a "Look" command at this scene will instead print the Crime Scene scenario.
    user.location = 2;
  }

  //Crime Scene
  if loc == 2 {
    fmt.Print("You stand at the scene of the crime. Blood covers the front entrance to Michelson. There appear to be some strange 'markings' in the" +
    " blood. Mr. Holmes waits off to the side, next to the body of the victim.\n\n");
  }

  //Ground Floor
  if loc == 3 {
    fmt.Printf("The hallways inside Michelson are deathly quiet. Before you is the staircase leading up to the \"Top Floor\" of Michelson. There" +
    " is nothing here except for a small streak of blood leading up the stairs.\n\n");
  }

  //Top Floor
  if loc == 4 {
    fmt.Printf("The upper deck of Michelson is deathly cold, but not out of the ordinary. To your left are the lab rooms (Lab 1, Lab 2, and Lab 3)" +
    " and to your right is Dr. Aviv's office, Dr. Roche's office, and Capt Bilzor's office.\n\n");
  }

  //Lab 1
  if loc == 5 {
    fmt.Printf("You hear the quiet humm of the lab machines. Although all the lights are off, you can make out the outlines of a few overturned lab" +
    " chairs. By the chairs, you see a small lump but are unsure what it may be.\n\n");
  }

  //Lab 2
  if loc == 6 {

    //Checks that user found the cipher text and has visited Lab 2 once before
    if (user.notepad[6] == true) && (user.visited[3] == true) {
      fmt.Printf("As you walk into the lab room, you notice that a workstation is still currently logged into. However, there are no signs that anyone" +
      " is currently using the room. You now notice that the closet door in the back of the room is slightly ajar.\n\n");
    } else {
      fmt.Printf("As you walk into the lab room, you notice that a workstation is still currently logged into. However, there are no signs that anyone" +
      " is currently using the room. A small closet takes up space in the back of the room, but the door is closed and the lights off.\n\n");
    }

    user.visited[3] = true;
  }

  //Lab 3
  if loc == 7 {
    fmt.Printf("All the lights are on, but there are no signs that anyone is currently occupying the room. You look around and notice a desk on the other" +
    " end of the room.\n\n");
  }

  //Dr. Roche's Office
  if loc == 8 {
    fmt.Printf("Dr. Roche's office is relatively bare except for a large desk that takes up a significant portion of the office space. A whiteboard hangs" +
    " on one of the walls.");

    //If this is user's first visit, then Dr. Roche will be in his Office. Otherwise check the randomly set variable for availability
    if (user.visited[1] == false) || (present == 1) {
      fmt.Printf(" Dr. Roche is sitting at his desk working diligently on something.\n\n");
    } else {
      fmt.Printf("\n\n");
    }
  }

  //Dr. Aviv's Office
  if loc == 9 {
    fmt.Printf("Dr. Aviv's office is cluttered with all sorts of knick-knacks. His desk has a multitude of material strewn about and his book shelf is" +
    " overflowing with books. He has a small closet tucked into the corner of his office.");

    //If this is user's first visit, then Dr. Aviv will be in his Office. Otherwise check the randomly set variable for availability
    if (user.visited[0] == false) || (present == 0) {
      fmt.Printf(" Dr. Aviv is currently in his office, browsing the internet, probably watching cat videos.\n\n");
    } else {
      fmt.Printf("\n\n");
    }
  }

  //Capt Bilzor's Office
  if loc == 10 {
    fmt.Printf("Capt Bilzor's office is in pristine condition. He is rapidly typing away at his computer and appears to be in deep concentration.\n\n");

    //Don't need to distinguish between first and subsequent visits unlike with Dr. Roche and Dr. Aviv
    user.visited[2] = true;
  }

  //Closet
  if loc == 11 {
    fmt.Printf("In the closet is an empty desk and a trashcan. You notice a small box stashed away in the darkest corner of the closet.\n\n");
  }
}


/*
 * Reads in the first word, known as the command, and calls the appropriate function
 */
func execute() {

  var command string;
  readcmd := true; //Used to determine if command needs to break from loop

  for readcmd {

    fmt.Printf("> ");
    fmt.Scanf("%s", &command); //Scanf only reads until a whitespace character
    command = strings.ToLower(command);

    //CallClear();

    //Call the appropriate function based on the command entered
    switch command {
      case "move":
        readcmd = move();

      case "examine":
        examine();

      case "ask" :
        ask();

      case "accuse":
        accuse();
        readcmd = false;

      case "look":
        fmt.Println("\n");
        prompt();

      case "notepad":
        notepad();

      case "save":
        save();

      case "quit":
        quit();
        readcmd = false;//stop reading commands after we quit

      case "help":
        help();

      default:/*DOESNT QUITE WORK RIGHT*/ //FIXME
        //Need to clear stdin buffer incase user typed additional words
        //https://coderwall.com/p/zyxyeg/golang-having-fun-with-os-stdin-and-shell-pipes
        fi,_ := os.Stdin.Stat()
        if fi.Size() > 0 {
          args := getArgs();
          //https://stackoverflow.com/questions/21743841/how-to-avoid-annoying-error-declared-and-not-used
          _ = args; //writes the string to an empty variable
        }

        fmt.Println("Unknown command");
        fmt.Println("Type \"Help\" for a list of commands");
    }
  }
}


/*
 * Moves the user to a specific location
 */
func move() bool {

  args := getArgs(); //The new location

  if (user.location == 0) && (strings.Compare(args, "crime scene") != 0){
    //Ensures that the first command execute is "Move Crime Scene"
    fmt.Printf("You must hurry to the \"Crime Scene\"\n\n");
    return true;
  } else if (user.location == 0) && (strings.Compare(args, "crime scene") == 0){
    //Moves user to the 2nd part of the intro
    user.location = 1;
    return false;
  } else {

    switch args{
      case "crime scene":
        user.location = 2;

      case "ground floor":
        user.location = 3;

      case "top floor":
        user.location = 4;

      case "lab 1":
        user.location = 5;

      case "lab 2":
        user.location = 6;

      case "lab 3":
        user.location = 7;

      case "dr. roche's office":
        user.location = 8;
        //Checks if it's user's first visit
        if user.visited[1] == false {
          user.visited[1] = true;
          present = 1; //Forces Dr. Roche to be in his office during first visit
        } else {
          present = rand.Intn(10) % 2; //Randomly determine if Dr. Roche will be in his office
        }

      case "dr. aviv's office":
        user.location = 9;
        //Checks if it's user's first visit
        if user.visited[0] == false{
          user.visited[0] = true;
          present = 0; //Forces Dr. Aviv to be in his office during first visit
        } else {
          present = rand.Intn(10) % 2; //Randomly determine if Dr. Aviv will be in his office
        }

      case "capt bilzor's office":
        user.location = 10;

      case "closet":
        //Only allow user to move to this location if they have visited Lab 2 before and have
        //found the cipher text
        if (user.notepad[6] == true) && (user.visited[3] == true) {
          user.location = 11;
        } else {
          fmt.Println("You cannot move into the closet.\n");
        }

      default:
        fmt.Printf("You cannot move to %s\n\n", args);
        //Causes the while loop in execute() to continue
        return true;
    }

    //Breaks out of the while loop in execute()  
    return false;
  }
}


/**
 * Determines the user's location and prints out the appropriate response depending on the item
 * the user is trying to examine. Also set the corresponding section in the user's notepad to true
 */
func examine() {

  fmt.Println("\n");
  args := getArgs();    //The item to be examine
  loc := user.location; //Current location

  switch loc{
    case 2:
      crimeScene(args);

    case 5:
      lab1(args);

    case 6:
      lab2(args);

    case 7:
      lab3(args);

    case 8:
      roche(args);

    case 9:
      aviv(args);

    case 11:
      closet(args);

    default:
      fmt.Printf("You cannot examine %s!\n\n", args);
  }
}


/*
 * Examine Crime Scene
 */
func crimeScene(args string) {

  if strings.Compare(args, "markings") == 0 {
    fmt.Println("You take a closer look at the strange markings. They appear to be footprints, approximately fitting a size 10 male shoe." +
    " They lead into the \"Ground Floor\" of Michelson.\n");
    user.notepad[0] = true;

  } else if strings.Compare(args, "victim") == 0 {
      fmt.Println("You take a closer look at the victim and recognize your former Algorithms Professor, Dr. Brown. His hands look pretty" +
      " beat up, like he put up a fight before his life ended.\n");
      user.notepad[1] = true;

    } else {
      fmt.Printf("You cannot examine %s!\n\n", args);
  }
}


/*
 * Examine Lab 1
 */
func lab1(args string) {

  if strings.Compare(args, "lump") == 0 {
    fmt.Println("It is a shoe! And it matches the size and tread of the shoe markings at the scene of the crime. The shoe is a black Nike.\n");
    user.notepad[2] = true;

  } else {
      fmt.Printf("You cannot examine %s!\n\n", args);
  }
}


/*
 * Examine Lab 2
 */
func lab2(args string) {

  if strings.Compare(args, "closet") == 0 {

    fmt.Println("The closet is locked. Would you like to try unlocking it? (Y/N)");

    answer := getArgs(); //User's response

    if strings.Compare(answer,"y") == 0 {
      fmt.Printf("Enter combination: ");
      lock := getArgs();
      _ = lock; //The user can never open the closet
      fmt.Println("Nothing happened...\n");

    } else if (strings.Compare(answer, "n") == 0) {
        //Do nothing if user types n

    } else {
        fmt.Printf("Please specifiy Y or N\n");
    }

  } else if strings.Compare(args, "workstation") == 0 {

      fmt.Println("Dr. Aviv is currently logged on. Would you like to try entering a password? (Y/N)");

      answer := getArgs(); //User's response

      if (strings.Compare(answer,"y") == 0) {
        fmt.Printf("Enter combination: ");
        var password string;
        fmt.Scanf("%s", &password);

        //Check is password is correct
        if strings.Compare(password, "c4tzWillrulD4W0rld!") == 0 {
          fmt.Println("Success! You notice that Dr. Aviv has his e-mail account pulled up. In his inbox are several angry e-mails addressed" +
          " to Dr. Brown dated merely hours before the indicated T.O.D. As you skim through the e-mails, you note that in one of Aviv's replies," +
          " he states that Dr. Brown \"will pay for his transgressions\".\n");
          user.notepad[4] = true;

        } else {
            fmt.Println("Nothing happened...\n");
        }

      } else if (strings.Compare(answer, "n") == 0) {
          //Do nothing if user types n
      } else {
        fmt.Printf("Please specifiy Y or N\n");
      }

  } else {
      fmt.Printf("You cannot examine %s!\n\n", args);
  }
}


/*
 * Examine Lab 3
 */
func lab3(args string) {

  if strings.Compare(args, "desk") == 0 {
    fmt.Println("On top of the desk is a single folded sheet of paper.\n");

  } else if strings.Compare(args, "paper") == 0 {
      fmt.Println("You unfold the paper and see a bunch of cat drawings all over it. There is also an alphanumeric sequence circled" +
      "several times: \"c4tzWillrulD4W0rld!\".\n");
      user.notepad[5] = true;

  } else {
      fmt.Printf("You cannot examine %s!\n\n", args);
  }
}


/*
 * Examine Dr. Roche's Office
 */
func roche(args string) {

  //User cannot examine objects in the room when the professor is present
  if present == 1 {
    if strings.Compare(args, "dr. roche") == 0 {
      fmt.Println("Nothing looks out of the ordinary with Dr. Roche's appearance.\n");
      user.notepad[10] = true;

    } else {
        fmt.Printf("You cannot examine %s!\n\n", args);
    }

  } else {
      if strings.Compare(args, "desk") == 0 {
        fmt.Println("You scan Dr. Roche's desk looking for any evidence. You see a stack of papers and a book on Algorithms.\n");

      } else if strings.Compare(args, "papers") == 0 {
          fmt.Println("Appears to be a stack of student HW that is not yet graded.\n");

      } else if strings.Compare(args, "book") == 0 {
          fmt.Println("You open the book and on the inside is a scrap of paper. The paper has \"PQXOQPUE\" written on it. What could" +
          " it possibly mean?\n");
          user.notepad[6] = true;

      } else if strings.Compare(args, "whiteboard") == 0 {
          fmt.Println("There is a lot of writing on the board but you cannot make sense of any of it.\n");

      } else {
          fmt.Printf("You cannot examine %s!\n\n", args);
      }
  }
}


/*
 * Examine Dr. Aviv's Office
 */
func aviv(args string) {

  //User cannot examine objects in the room when the professor is present
  if present == 0 {
    if strings.Compare(args, "dr. aviv") == 0 {
      fmt.Println("You notice that Dr. Aviv appears a bit haggard. He also has several scratch marks on his hands.\n");
      user.notepad[9] = true;

    } else {
      fmt.Printf("You cannot examine %s!\n\n", args);
    }

  } else {
    if strings.Compare(args, "desk") == 0 {
      fmt.Println("There are so many papers covering his desk you aren't sure where to start. None of the papers look particularly interesting.\n");

    } else if strings.Compare(args, "book shelf") == 0 {
        fmt.Println("As you flip through some of the books, you do not find anything of evidentiary value. However, you do notice a sticky note," +
        "labeled \"Locker Combos\" hanging off the shelf.\n");

    } else if strings.Compare(args, "sticky note") == 0 {
        fmt.Println("There are several combinations written on the note: \"23-4-12-6\", \"56-8-3-0\", \"48-93-2-43\"\n");
        user.notepad[7] = true;

    } else if strings.Compare(args, "closet") == 0 {
        fmt.Println("You notice a slight stench emanating from the closet. You slowly open the door, a slight squeak breaking the silence in the room." +
        "At the bottom of the closet, you notice a shoe, missing its other pair. The shoe is a size 10 black Nike.\n");
        user.notepad[8] = true;

    } else {
        fmt.Printf("You cannot examine %s!\n\n", args);
    }
  }
}


/*
 * Examine Closet
 */
func closet(args string) {

  if strings.Compare(args, "trashcan") == 0 {
    fmt.Println("There are some discarded coffee cups and a crumbled up piece of paper.\n");

  //Cipher alphabet needed to decode the cipher text
  } else if strings.Compare(args, "paper") == 0 {
      fmt.Println("On the paper is written: CIPHERFUNAJLZVQBXOYSGWMDKT.\n");
      user.notepad[3] = true;

  } else if strings.Compare(args, "box") == 0 {
      fmt.Println("It is locked, but there is a keypad available to enter a passcode of some sort. Would you like to try unlocking it? (Y/N)");

      answer := getArgs();

      if (strings.Compare(answer,"y") == 0) {
        fmt.Printf("Enter passcode: ");
        password := getArgs();

        if strings.Compare(password, "coqroche") == 0 {
          fmt.Println("");
          fmt.Println("You hear a click as the lip of the box pops open. Inside the box, you find several papers written by Dr. Brown proving that" +
          " P = NP. You realize that if Dr. Brown went public with this proof, he would have been instantly rich and famous! As you continue to" +
          " inspect the papers, a photograph falls to the floor. It is a photograph of Dr. Brown with a red X drawn through his face. Looks like you" +
          " found the killer.\n");
          user.notepad[11] = true;
        } else {
            fmt.Println("Nothing happened...\n");
        }
      } else if (strings.Compare(answer, "n") == 0) {

      } else {
          fmt.Printf("Please specifiy Y or N\n");
      }

  } else {
      fmt.Printf("You cannot examine %s!\n\n", args);
  }
}


/*
 * Ask a given suspect about one of the other professors. 
 */
func ask() {

  fmt.Println("\n");
  loc := user.location; //Current location

  switch loc {
    case 8:
      //A professor has to be in his office in order to question him
      if present == 1 {
        question("Dr. Roche", "Dr. Aviv", "Capt Bilzor");
      } else {
        fmt.Println("Dr. Roche has to be present to ask him any questions!");
      }
    case 9:
      if present == 0 {
        question("Dr. Aviv", "Dr. Roche", "Capt Bilzor");
      } else {
        fmt.Println("Dr. Aviv has to be present to ask him any questions!");
      }
    case 10:
      question("Capt Bilzor", "Dr. Aviv", "Dr. Roche");
    default:
      fmt.Println("You have to be in an office to ask a question!");
  }
}


/*
 * The user chooses a question from a list and then the appropriate response is generated
 */
func question(suspect string, other1 string, other2 string) {

  //Keep looping until user inputs an appropriate response
  for true {
    fmt.Println("\n");
    fmt.Printf("1) Ask %s where he was at the time of the murder.\n", suspect);
    fmt.Printf("2) Ask %s about %s.\n", suspect, other1);
    fmt.Printf("3) Ask %s about %s.\n", suspect, other2);
    fmt.Printf("4) Ask %s when he last saw Dr. Brown.\n", suspect);
    fmt.Println("Which question would you like to ask (1/2/3/4): ");

    var answer int;
    fmt.Scanf("%d", &answer);

    if answer == 1 {

      if strings.Compare(suspect, "Dr. Roche") == 0{
          fmt.Printf("Did you forget already, %s! I briefly stopped by for Pints with Profs the other night. I said hi to a few people before" +
          " leaving. You looked like you were having a lot of fun!\n\n", user.name);

        } else if strings.Compare(suspect, "Dr. Aviv") == 0{
            fmt.Println("I left early to pick up the new cat I adopted. I've been taking care of her all weekend! It's a tiring job.\n");

        } else {
            fmt.Println("Although I've been really busy with work here, I did stop by Pints with Profs for a little bit.\n");
        }
        //Printed the appropriate response and break out of loop
        break;

    } else if answer == 2 {

        if strings.Compare(suspect, "Dr. Roche") == 0{
          fmt.Println("He wasn't at Pints with Profs so I have no idea what he was up to the other night.\n");

        } else if strings.Compare(suspect, "Dr. Aviv") == 0{
            fmt.Println("I have no idea what Dr. Roche was up to the other night. I was at home watching cat videos.");

        } else {
          fmt.Println("Dr. Aviv left early that day to go pick up a cat he just recently adopted. He couldn't stop talking about it all day!" +
          " I didn't see him again after that, but I'm assuming he was at home taking care of his new pet.\n");
        }
        //Printed the appropriate response and break out of loop
        break;

    } else if answer == 3 {

        if strings.Compare(suspect, "Dr. Roche") == 0{
          fmt.Println("The Computer Science Department was swamped with a lot of tasks to complete this last week. I know Capt Bilzor had been" +
          " working tirelessly to make sure everything was complete. I wouldn't be surprised Vif he stayed late in order to get everything done.\n");

        } else if strings.Compare(suspect, "Dr. Aviv") == 0{
            fmt.Println("Capt Bilzor has been busy with Department stuff all week. I wouldn't be surprised if he just stayed here over the weekend" +
            " to finished it!\n");

        } else {
          fmt.Println("I remember Dr. Roche being excited about something before he left for the day. I'm not sure what he was excited about though." +
          " However, he did say he would be busy and wouldn't be able to stop by Pints with Profs. I didn't see him there either.\n");
        }
        //Printed the appropriate response and break out of loop
        break;

    } else if answer == 4 {

        if strings.Compare(suspect, "Dr. Roche") == 0{
          fmt.Println("The last time I saw Dr. Brown was on Friday before I departed work for the day. He seemed to be really excited about something" +
          " and said it would change everything we knew about Computer Science. I have no idea what he could have been referring to though.\n");

        } else if strings.Compare(suspect, "Dr. Aviv") == 0{
            fmt.Println("Dr. Brown and I had a disagreement this last week so I hadn't seen him.\n");

        } else {
          fmt.Println("The last time I saw Dr. Brown was on Friday when he departed work for the day. He was very excitedly discussing something with" +
          " Dr. Roche before he left.");
        }
        //Printed the appropriate response and break out of loop
        break;

    } else {
        fmt.Printf("Please specifiy 1, 2, 3, or 4\n");
    }
  }
}


/*
 * Accuse an individual of being the killer!
 */
func accuse() {

  fmt.Println("\n");
  args := getArgs();    //Who the user is accusing
  loc := user.location; //Current location

  //Detective Holmes said we had to talk to him inorder to accuse someone!  
  if loc != 2 {
    fmt.Println("You have to be at the \"Crime Scene\n in order to accuse someone!\n");

  } else {

      if strings.Compare(args, "dr. aviv") == 0{

        if (user.notepad[2] == true) && (user.notepad[8] == true) {
          fmt.Println("You submit the evidence to Detective Holmes and the police swiftly enter Michelson to apprehend the suspect. That night, you" +
          " climb into bed feeling proud that you were able to catch the murderer. You fall into a deep, dreamless sleep, none the wiser of the" +
          " mysterious shadow that stands just outside your window. You awake, heart pounding, as someone whispers into your ear.\n\n \"You got the" +
          " wrong guy.\"\n\n That is the last thing you hear before everything goes completely dark...\n\n");
          quit();

        } else {
            fmt.Println("You do not have enough evidence to accuse this individual.\n");
        }

      } else if strings.Compare(args, "dr. roche") == 0 {

          if user.notepad[11] == true {
            fmt.Println("You submit the evidence to Detective Holmes and the police swiftly enter Michelson to apprehend the suspect. Congratulations!" +
            " You found the killer. Dr. Roche was jealous of Dr. Brown's impending fame and fortune, so Dr. Roche conspired to get rid of Dr. Brown so" +
            " that he could claim the proof as his own.\n\n Computer Science is risky buisness.\n\n");
            quit();

          } else {
            fmt.Println("You do not have enough evidence to accuse this individual.\n");
          }

      } else if strings.Compare(args, "capt bilzor") == 0 {
          fmt.Println("You do not have enough evidence to accuse this individual.\n");

      } else {
          fmt.Printf("You cannot accuse %s!\n\n", args);
      }
  }
}


/*
 * Prints out a list of all the evidence the user has found
 */
func notepad() {

  fmt.Println("\n");

  notes := [12]string{
      "Bloddy footprints in a male size 10 shoe.",
      "Victim looks like he put up a fight.",
      "Shoe found in Lab 1 that matches the footprints from the crime scene. A black Nike.",
      "CIPHERFUNAJLZVQBXOYSGWMDKT found in Dr. Roche's office.",
      "Aviv sent several angry e-mails to Dr. Brown just hours before he died.",
      "\"c4tzWillrulD4W0rld!\" found on a pice of paper with cats drawn all over it.",
      "\"PQXOQPUE\" found in Dr. Roche's office. What does it mean?",
      "Several combinations written on a note in Dr. Aviv's Office: \"23-4-12-6\", \"56-8-3-0\", \"48-93-2-43\"",
      "A black Nike shoe, size 10, from Dr. Aviv's office.",
      "Dr. Aviv looks haggard and has various scartch marks all over his hands.",
      "Nothing out of the ordinary with Dr. Roche's appearance.",
      "A box containing incriminating evidence." }

  //Checks which evidence in the notepad is true then prints out the associated evidence
  for i := range user.notepad {
    if user.notepad[i] == true {
      fmt.Printf("%s\n", notes[i]);
    }
  }

  fmt.Println("");
}

func save() {
  fd, _ := os.OpenFile("./.Murder", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666);
  writer := bufio.NewWriter(fd) //http://codefoolz.com/blog/golang/2015/07/25/Writing-to-File-in-Go-lang.html
  defer fd.Close();

  fmt.Fprintf(writer, "%s\n", user.name)
  fmt.Fprintf(writer, "%d\n", user.location)
  for i:=0;i<11;i++{
    if(user.notepad[i] == true){ //printing a bool in the format that I wanted was troublesome
      fmt.Fprintf(writer, "%s,", "true");
    }else{
      fmt.Fprintf(writer, "%s,", "false");
    }
  }
  if(user.notepad[11] == true){ //printing a bool in the format that I wanted was troublesome
      fmt.Fprintf(writer, "%s\n", "true");
    }else{
      fmt.Fprintf(writer, "%s\n", "false");
    }
  for i:=0;i<3;i++{
    if(user.visited[i] == true){ //printing a bool in the format that I wanted was troublesome
      fmt.Fprintf(writer, "%s,", "true");
    }else{
      fmt.Fprintf(writer, "%s,", "false");
    }
  }
  if(user.visited[3] == true){ //printing a bool in the format that I wanted was troublesome
      fmt.Fprintf(writer, "%s\n", "true");
    }else{
      fmt.Fprintf(writer, "%s\n", "false");
    }
  writer.Flush()

}

func quit() {
  fmt.Println("");
  //channel to send the exit to other threads
  quitting = 1; //<- true;

}

func help() {
  fmt.Println("\n");
  fmt.Println("The murderer must be found before he claims more victims. Investigate, question, and discover the truth.\n");
  fmt.Println("Your command words are:");
  fmt.Println("Move, Examine, Ask, Accuse, Look, Notepad, Save, Quit\n");
  fmt.Println("To speak to potential suspects you need to ask them something. Type \"Ask\" while in an office and the following list of" +
  "questions will be made available:");
  fmt.Println("1) Ask <name> where they were at the time of the murder.");
  fmt.Println("2) Ask <name> about <name>.");
  fmt.Println("3) Ask <name> about <name>.");
  fmt.Println("4) Ask <name> about when they last saw the victim alive.\n");
  fmt.Println("Any clues you find will be automatically added to your notepad. You can view all your clues by typing \"Notepad\".\n");
}


/*
 * Reads any arguements that come after the command word
 */
func getArgs() string {
  //https://stackoverflow.com/questions/14094190/golang-function-similar-to-getchar
  //Reads up to aand including '\n'
  reader := bufio.NewReader(os.Stdin);
  args, _ := reader.ReadString('\n');
  args = strings.ToLower(args); //everything is lowercase to simplify comparisons
  args = strings.TrimSuffix(args, "\n");
  return args;
}


/*
 * Initiates the sequence to play MICHELSON MURDER MYSTERY!
 */
func game(wg * sync.WaitGroup) {
  defer wg.Done();

  //CallClear();
  fmt.Println();
  fmt.Printf(
` ███▄ ▄███▓ ██▓ ▄████▄   ██░ ██ ▓█████  ██▓      ██████  ▒█████   ███▄    █    
▓██▒▀█▀ ██▒▓██▒▒██▀ ▀█  ▓██░ ██▒▓█   ▀ ▓██▒    ▒██    ▒ ▒██▒  ██▒ ██ ▀█   █    
▓██    ▓██░▒██▒▒▓█    ▄ ▒██▀▀██░▒███   ▒██░    ░ ▓██▄   ▒██░  ██▒▓██  ▀█ ██▒   
▒██    ▒██ ░██░▒▓▓▄ ▄██▒░▓█ ░██ ▒▓█  ▄ ▒██░      ▒   ██▒▒██   ██░▓██▒  ▐▌██▒   
▒██▒   ░██▒░██░▒ ▓███▀ ░░▓█▒░██▓░▒████▒░██████▒▒██████▒▒░ ████▓▒░▒██░   ▓██░   
░ ▒░   ░  ░░▓  ░ ░▒ ▒  ░ ▒ ░░▒░▒░░ ▒░ ░░ ▒░▓  ░▒ ▒▓▒ ▒ ░░ ▒░▒░▒░ ░ ▒░   ▒ ▒    
░  ░      ░ ▒ ░  ░  ▒    ▒ ░▒░ ░ ░ ░  ░░ ░ ▒  ░░ ░▒  ░ ░  ░ ▒ ▒░ ░ ░░   ░ ▒░   
░      ░    ▒ ░░         ░  ░░ ░   ░     ░ ░   ░  ░  ░  ░ ░ ░ ▒     ░   ░ ░    
       ░    ░  ░ ░       ░  ░  ░   ░  ░    ░  ░      ░      ░ ░           ░    
               ░                                                               
          ███▄ ▄███▓ █    ██  ██▀███  ▓█████▄ ▓█████  ██▀███                   
         ▓██▒▀█▀ ██▒ ██  ▓██▒▓██ ▒ ██▒▒██▀ ██▌▓█   ▀ ▓██ ▒ ██▒                 
         ▓██    ▓██░▓██  ▒██░▓██ ░▄█ ▒░██   █▌▒███   ▓██ ░▄█ ▒                 
         ▒██    ▒██ ▓▓█  ░██░▒██▀▀█▄  ░▓█▄   ▌▒▓█  ▄ ▒██▀▀█▄                   
         ▒██▒   ░██▒▒▒█████▓ ░██▓ ▒██▒░▒████▓ ░▒████▒░██▓ ▒██▒                 
         ░ ▒░   ░  ░░▒▓▒ ▒ ▒ ░ ▒▓ ░▒▓░ ▒▒▓  ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░                 
         ░  ░      ░░░▒░ ░ ░   ░▒ ░ ▒░ ░ ▒  ▒  ░ ░  ░  ░▒ ░ ▒░                 
         ░      ░    ░░░ ░ ░   ░░   ░  ░ ░  ░    ░     ░░   ░                  
                ░      ░        ░        ░       ░  ░   ░                      
                                       ░                                       
       ███▄ ▄███▓▓██   ██▓  ██████ ▄▄▄█████▓▓█████  ██▀███ ▓██   ██▓           
      ▓██▒▀█▀ ██▒ ▒██  ██▒▒██    ▒ ▓  ██▒ ▓▒▓█   ▀ ▓██ ▒ ██▒▒██  ██▒           
      ▓██    ▓██░  ▒██ ██░░ ▓██▄   ▒ ▓██░ ▒░▒███   ▓██ ░▄█ ▒ ▒██ ██░           
      ▒██    ▒██   ░ ▐██▓░  ▒   ██▒░ ▓██▓ ░ ▒▓█  ▄ ▒██▀▀█▄   ░ ▐██▓░           
      ▒██▒   ░██▒  ░ ██▒▓░▒██████▒▒  ▒██▒ ░ ░▒████▒░██▓ ▒██▒ ░ ██▒▓░           
      ░ ▒░   ░  ░   ██▒▒▒ ▒ ▒▓▒ ▒ ░  ▒ ░░   ░░ ▒░ ░░ ▒▓ ░▒▓░  ██▒▒▒            
      ░  ░      ░ ▓██ ░▒░ ░ ░▒  ░ ░    ░     ░ ░  ░  ░▒ ░ ▒░▓██ ░▒░            
      ░      ░    ▒ ▒ ░░  ░  ░  ░    ░         ░     ░░   ░ ▒ ▒ ░░             
             ░    ░ ░           ░              ░  ░   ░     ░ ░                
                  ░ ░                                       ░ ░             `);
  fmt.Printf("\n\n\t\tby Brandon Sipes and Kristina Bodeman\n\n\n\n");

  //Checks for a previously saved file
  if (checkSave()) {

    fmt.Printf("Would you like to load the saved file? (Y/N)\n");

    for true {
      answer := getArgs();

      if strings.Compare(answer,"y") == 0 {
        fmt.Printf("Loaded saved file\n\n");
        setUser();
        break;

      } else if strings.Compare(answer, "n") == 0 {
        fmt.Printf("\n\nEnter your name: ");
        user.name = getArgs();
        break;

      } else {
        fmt.Printf("Please specifiy Y or N\n");
      }
    }

  } else {
    fmt.Printf("Enter your name: ");
    user.name = getArgs();
  }

  //Continuously prints out the prompt for user's current location and then
  //executes the next command
  for quitting == 0{
    prompt();
    execute();
  }
  /*var run bool = true;
  for run{ //TODO make channel that quit() can change
    switch{
    case <-quitting:
      fmt.Println("TEST2");
      break;
    default:
      fmt.Println("TEST3");
      //CallClear();
      prompt();
      execute();
    }
  }*/
}

//https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go/22892171
//Golang does not have a built in clear screen function. Found this on stack overflow.
var clear map[string]func() //create a map for storing clear funcs

func init() {
  clear = make(map[string]func())
  clear["linux"] = func() {
    cmd := exec.Command("clear") //Linux example
    cmd.Stdout = os.Stdout
    cmd.Run()
  }
  clear["windows"] = func() {
    cmd := exec.Command("cmd", "/c", "cls") //Windows example
    cmd.Stdout = os.Stdout
    cmd.Run()
  }
}

func CallClear() {
  value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
  if ok { //if we defined a clear func for that platform:
    value()  //we execute it
  } else { //unsupported platform
      panic("Your platform is unsupported! I can't clear terminal screen :(")
  }
}
func dogecoin(){
  //You found a hash!!
  //Time to submit it to dogecoin database to get
  //more dogecoins!!

  //Apparently it is against ITSD rules to mine crypto currency...
  //So this function does not connact the dogecoin database 
}

func doge(wg * sync.WaitGroup){
  defer wg.Done();

  var hash []byte;
  if _, err := os.Stat("./.hash"); !os.IsNotExist(err){
    fd,_ := os.Open("./.hash");

    scanner := bufio.NewScanner(fd);//create a basic scanner
    scanner.Scan();
    hash = scanner.Bytes();//read line in bytes

    fd.Close();
  }else{ //No file exists
    data := make([]byte, 12);
    rand.Read(data); //load it up with random bytes

    sum := sha256.Sum256(data);
    hash = sum[:]; //splice
  }

  var loop int;
  for loop == 0 {//TODO end when quit() is called
    select{
    case  <- quitting:
      loop = 1;
    default:
      sum := sha256.Sum256(hash);
      hash = sum[:];
      if hash[0] == 0 && hash[1] == 0{
        dogecoin();
      }
    }
  }
  //for quitting == 0{
  //  sum := sha256.Sum256(hash);
  //  hash = sum[:];
  //  if hash[0] == 0 && hash[1] == 0{
  //    dogecoin();
  //  }
  //}

  hashfd,_ := os.OpenFile("./.hash", os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0666);
  writer := bufio.NewWriter(hashfd);
  defer hashfd.Close();

  fmt.Fprintf(writer, "%d", hash);
  writer.Flush()//save hash
}

func main() {

  //Adds the time as a seed for the random function
  rand.Seed(time.Now().UTC().UnixNano());
  sigc := make(chan os.Signal, 1)//signal handling
  signal.Notify(sigc,
    syscall.SIGINT,
    syscall.SIGTERM)
  go func() {
    <-sigc
    quit()
  }()

  var wg sync.WaitGroup //https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html
  //Provides synchronization functions for go routines (threads).
  //Indicates that there will be 2 threads that must be waited on. Sets the WaitGroup counter to 2.
  wg.Add(1)
  //Executes the game as a separate thread
  go game(&wg)

  wg.Add(1)
  go doge(&wg)

  //Waits until the WaitGroup counter is 0.
  wg.Wait()
}
