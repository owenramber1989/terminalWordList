# terminalWordList

I'm a novice go programmer. This CLI wordlist tool was invented to help me better remember the new words or new phrases that I've encounterd. And it was built upon mysql.

## Features

1. Comprehensive help infomation

<img width="892" alt="image" src="https://user-images.githubusercontent.com/97386910/235935225-f652603f-fbd3-4c30-9770-d0c9aa8b88c8.png">


2. Complete vocabulary memorization process

We offered the insert, show, update importance, query, delete method for you. With these functions you shall be able to learn the new words easily!

- insert&update

<img width="1512" alt="image" src="https://user-images.githubusercontent.com/97386910/235547807-6eca858b-9cb4-4f01-85c8-f3fdfb38f7a3.png">

- delete

<img width="1307" alt="image" src="https://user-images.githubusercontent.com/97386910/235547930-2a9031f2-6dbc-4bee-8a55-d67e64f7145d.png">

- persistency

<img width="1512" alt="image" src="https://user-images.githubusercontent.com/97386910/235547977-3fe45f9a-8b40-42e7-8805-f19a1ee56d24.png">

3. Ebbinghaus Learing Curve

We use cycle to identify your review times and state to indentify your current learning phase. There are nine phases in all.

  1) 5m
  2) 30m
  3) 12h
  4) 1d
  5) 2d
  6) 4d 
  7) 7d
  8) 15d
  9) 31d
  
When you choose the operation -r, you reviewed your words/phrases. So the cycle will be added 1. And we measure the time passed since you first add this word/phrase into your database, which helps us find your supposed learning phase. If your current cycle is ahead of your current state AND your supposed phase if ahead of your current state, your current state will be added 1, too.

<img width="1210" alt="image" src="https://user-images.githubusercontent.com/97386910/235933176-137bba6a-1b38-4e40-b1fb-6c9235e18374.png">

4. recall the last review

<img width="1289" alt="image" src="https://user-images.githubusercontent.com/97386910/235934919-82906690-d83e-4646-bbaf-14e579dd134e.png">




## Get started!

```zsh
git clone https://github.com/owenramber1989/terminalWordList.git;
cd terminalWordList;
make;
sudo cp ./voc /usr/local/bin;
md ~/.config/voc/;
sudo cp ./voc.conf ~/.config/voc/;
voc -h;
```

## Configure

You can configure your database info in the file voc.conf. The default options is

- user: root
- password:
- host: localhost
- database: vocabulary

## Usage

1. Get the help message

`./voc -h`
