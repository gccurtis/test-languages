package main

import(
	"fmt"
	"strings"
	"strconv"
	"os"
	"bufio"
)

func abs(num int) (int){
	if num > 0{
		return num;
	} else{
		return -num;
	}
}

func handle(err error) (){
	if err != nil{
		panic(err);
	}
}

//Finds the matching parenthesis even if there parenthesis inbetween
func closingIndex(list []string, start int) (int){
	check := 1;
	pos := start;
	for check != 0{
		pos = pos + 1;
		if list[pos] == "("{
			check = check + 1;
		} else if list[pos] == ")"{
			check = check - 1;
		} else{
			continue;
		}
	}
	return pos;
}

//separates a list by parenthesis elements
func deparenth(list []string) ([][]string){
	var ret [][]string = make([][]string,0);
	pos := 0;
	for pos < len(list){
		if list[pos] == "("{
			end := closingIndex(list,pos);
			ret = append(ret,list[pos+1:end]);
			pos = end;
		} else{
			ret = append(ret,[]string{list[pos]});
		}
		pos = pos + 1;
	}
	return ret
}

//puts a list of list of strings together with parenthesis as separators
func listMerge(lili [][]string) ([]string){
	fmt.Println("a");
	fmt.Println(lili);
	var ret []string = make([]string,0);
	for _,inner_list := range lili{
		ret = append(append(ret,"("),append(inner_list,")")...);
	}
	fmt.Println("z");
	fmt.Println(ret);
	return ret;
}

func listCombine(lili []string) (string){
	var ret string = "";
	for _,char := range lili{
		ret = ret+char;
	}
	return ret;
}

func listCombinef(lili []string) (string){
	if len(lili) == 0{
		return "";
	}
	var ret string = "";
	for _,char := range lili[1:]{
		ret = ret+" "+char;
	}
	return lili[0]+ret;
}

func superlistCombine(lili [][]string) ([]string){
	var ret []string = make([]string,0);
	for _,li := range lili{
		ret = append(ret,li...);
	}
	return ret;
}

func listSplit(list []string, match string) ([][]string){
	ret := make([][]string,0)
	tmp := make([]string, 0)
	for _,word := range list{
		if word == match{
			ret = append(ret,tmp)
			tmp = make([]string,0)
		} else{
			tmp = append(tmp,word)
		}
	}
	ret = append(ret,tmp)
	return ret
}

//turns a string into a list of tokens
func tokenize (text string) ([]string){

	var letters []rune = []rune(text);
	var word string = "";
	var ret []string = make([]string,0);
	var quoteMode bool = false;

	for _,letter := range letters{
		if letter == '"'{
			if quoteMode && word != ""{
				ret = append(ret,word);
				word = "";
			}
			quoteMode = !quoteMode;
		} else if !quoteMode && letter == ' ' && word != ""{
			ret = append(ret,word);
			word = "";
		} else if !quoteMode && (letter == '&' || letter == '`'  || letter == '(' || letter == ')' || letter == '~'){
			if word != "" && word != " "{
				ret = append(ret,word);
			}
			ret = append(ret,string(rune(letter)));
			word = "";
		} else{
			if quoteMode{
				word = word + string(rune(letter));
			} else if letter != ' '{
				word = word + string(rune(letter));
			}
		}
	}

	if word != "" && word != " "{
		ret = append(ret,word);
	}
	return ret;
}

type Environment struct{
	environments []map[string][]string
	forms []map[string][]string
}

func newEnvironment() (*Environment){
	var ret *Environment = new(Environment);
	first := make([]map[string][]string,0);
	first = append(first,make(map[string][]string));
	ret.environments = first;
	second := make([]map[string][]string,0);
	second = append(second,make(map[string][]string));
	ret.forms = second;
	return ret;
}

func (env Environment) get(name string) ([]string){
	for _,environment := range env.environments{
		if val,ok := environment[name]; ok{
			return val;
		}
	}
	panic(name+" not in environment");
	return []string{};
}

func (env Environment) check(name string) ([]string,bool){
	for _,environment := range env.environments{
		if val,ok := environment[name]; ok{
			return val,true;
		}
	}
	return []string{},false;
}

func (env *Environment) update(name string, val []string) (){
	env.environments[0][name] = val;
}

func (env *Environment) extend() (){
	env.environments = append([]map[string][]string{make(map[string][]string)},env.environments...);
	env.forms = append([]map[string][]string{make(map[string][]string)},env.forms...);
}

func (env *Environment) contract() (){
	env.environments = env.environments[1:];
	env.forms = env.forms[1:];
}

func (env *Environment) apply_forms(possible []string) ([]string, bool){
	var fullcode [][]string = deparenth(possible);
	var form_as_list []string;
	for _,forms := range env.forms{
		for form,app := range forms{
			form_as_list = strings.Split(form," ");
			//fmt.Println("FORMS");
			//fmt.Println(form_as_list);
			if matches_form(fullcode,form_as_list){
				//fmt.Println("LOOKING");
				//fmt.Println(fullcode);
				return apply_form(fullcode,form_as_list,app,env),true;
			}
		}
	}
	return possible,false;
}

func (env *Environment) set_forms(form string, expr []string) (){
	env.forms[0][form] = expr;
}

//This is how you start the REPL or read a file
func main(){
	var env *Environment = newEnvironment();
	env.update("+",[]string{"native","|","+","|","&l"});
	env.update("-",[]string{"native","|","-","|","&l"});
	env.update("*",[]string{"native","|","*","|","&l"});
	env.update("/",[]string{"native","|","/","|","&l"});
	env.update("abs",[]string{"native","|","abs","|","&l"});
	env.update("number?",[]string{"native","|","number?","|","&l"});
	env.update("symbol?",[]string{"native","|","symbol?","|","&l"});
	env.update("eq?",[]string{"native","|","eq?","|","&l"});
	env.update("<",[]string{"native","|","<","|","&l"});
	env.update(">",[]string{"native","|",">","|","&l"});
	env.update(">=",[]string{"native","|",">=","|","&l"});
	env.update("<=",[]string{"native","|","<=","|","&l"});
	env.update("print",[]string{"native","|","print","|","&l"});
	env.update("println",[]string{"native","|","println","&l"});

	 /*for{
	 	scanner := bufio.NewScanner(os.Stdin);
	 	fmt.Print("> ");
	 	scanner.Scan();
	 	eval_expr(tokenize(scanner.Text()),env);
	 }*/

	file,err := os.Open("teszt0");
	handle(err);

	reader := bufio.NewReader(file);
	for true{
		inp,_,err := reader.ReadLine();
		strinp := string(inp);
		handle(err);
		if strinp == "exit"{
			os.Exit(0);
		}else{
			toks := tokenize(strinp);
			if len(toks) != 0 && toks[0] != ";"{
				//fmt.Print("COPY: ");
				//fmt.Println(toks);
				//fmt.Print("ENV: ");
				//fmt.Println(env);
				eval_expr(toks,env);
			}
		}
	}
}

func free(code []string, env *Environment) ([]string){
	var ret []string = make([]string,0);
	for _,tok := range code{
		if _,ok := env.check(tok); ok{
			ret = append(ret,free(env.get(tok),env)...);
		} else{
			ret = append(ret,tok);
		}
	}
	return ret;
}

//For builtin functions
func native_eval(oper string, params [][]string) ([]string){
	if oper == "+"{
		ret := 0;
		var tmp int;
		var err error;
		for _,param := range params{
			tmp,err = strconv.Atoi(param[0]);
			handle(err);
			ret = ret+tmp;
		}
		return []string{strconv.Itoa(ret)};
	} else if oper == "*"{
		ret := 1;
		var tmp int;
		var err error;
		for _,param := range params{
			tmp,err = strconv.Atoi(param[0]);
			handle(err);
			ret = ret*tmp;
		}
		return []string{strconv.Itoa(ret)};
	} else if oper == "-"{
		ret,err := strconv.Atoi(params[0][0]);
		handle(err)
		var tmp int;
		if len(params) == 1{
			return []string{strconv.Itoa(-ret)};
		}
		for _,param := range params[1:]{
			tmp,err = strconv.Atoi(param[0]);
			handle(err);
			ret = ret-tmp;
		}
		return []string{strconv.Itoa(ret)};
	} else if oper == "/"{
		ret,err := strconv.Atoi(params[0][0]);
		handle(err)
		var tmp int;
		if len(params) == 1{
			return []string{strconv.Itoa(1/ret)};
		}
		for _,param := range params[1:]{
			tmp,err = strconv.Atoi(param[0]);
			handle(err);
			ret = ret/tmp;
		}
		return []string{strconv.Itoa(ret)};
	} else if oper == "abs"{
			var ret []string = make([]string,0);
			for _,param := range params{
				nparam,err := strconv.Atoi(param[0]);
				handle(err);
				ret = append(ret,strconv.Itoa(abs(nparam)));
			}
			return ret;
	} else if oper == "eq?" || oper == "="{
		var test string = params[0][0];
		ret := true;
		for _,param := range params[1:]{
			ret = ret && (test == param[0]);
			test = param[0];
		}
		if ret{
			return []string{"#t"};
		} else{
			return []string{"#f"};
		}
	} else if oper == "<"{
		var test string = params[0][0];
		ret := true;
		for _,param := range params[1:]{
			ret = ret && (test < param[0]);
			test = param[0];
		}
		if ret{
			return []string{"#t"};
		} else{
			return []string{"#f"};
		}
	} else if oper == ">"{
		var test string = params[0][0];
		ret := true;
		for _,param := range params[1:]{
			ret = ret && (test > param[0]);
			test = param[0];
		}
		if ret{
			return []string{"#t"};
		} else{
			return []string{"#f"};
		}
	} else if oper == "symbol?"{
		var ret bool = true;
		for _,param := range params{
			if _,err := strconv.Atoi(param[0]); err == nil{
				ret = false;
			}
		}
		if ret{
			return []string{"#t"};
		} else{
			return []string{"#f"};
		}
	} else if oper == "number?"{
		var ret bool = true;
		for _,param := range params{
			if _,err := strconv.Atoi(param[0]); err != nil{
				ret = false;
			}
		}
		if ret{
			return []string{"#t"}
		} else{
			return []string{"#f"}
		}
	}else{
		fmt.Println(oper+" is not an implemented builtin function");
		return []string{};
	}
}

func apply(proc []string, params [][]string, env *Environment) ([]string){
	tmp := listSplit(proc,"|");
	var check []string = tmp[0];
	if check[0] == "native"{
		return native_eval(tmp[1][0],params);
	}
	var args []string = tmp[1];
	var expr []string = tmp[2];
	env.extend();
	for ind,varname := range args{
		if(string(varname[0]) == "&"){
			env.update(args[ind+1],superlistCombine(params[ind:]));
			break;
		} else{
			env.update(varname,params[ind])
		}
	}
	ret := eval_expr(expr,env);
	env.contract();
	return ret;
}

func release_lambda(proc []string, params [][]string, env *Environment) ([]string){
	tmp := listSplit(proc,"|");
	if(len(tmp) != 3){
		return proc;
	}
	var args []string = tmp[1];
	var expr []string = tmp[2];
	env.extend();
	for ind,varname := range args{
		env.update(varname,params[ind]);
	}

	ret := eval_expr(append([]string{"rel"},expr...),env);
	env.contract();
	return ret;
}

func matches_form(fullcode [][]string, form []string) (bool){
	if(len(fullcode) == 0 || len(form) == 0){
		return false;
	}
	var i int = 0;
	for i < len(fullcode){
		if form[i] == "&"{
			return true;
		}
		if form[i] == "~" || form[i] == "`"{
			i = i+1;
		} else{
			if len(fullcode[i]) == 0{
				return false;
			}
			if form[i] != fullcode[i][0]{
				return false;
			}
		}
		i = i+1
	}
	return true;
}

func apply_form(fullcode [][]string, form []string, expr []string, env *Environment) ([]string){
	//var tech_env map[string][]string = make(map[string][]string);
	codepos := 0;
	formpos := 0;
	env.extend();
	for codepos < len(fullcode){
		if form[formpos] == "`"{
			env.update(form[formpos+1], fullcode[codepos]);
			formpos = formpos+1;
		} else if form[formpos] == "~"{
			env.update(form[formpos+1], eval_expr(fullcode[codepos],env));
			formpos = formpos+1
		} else if form[formpos] == "&"{
			env.update(form[formpos+1], superlistCombine(fullcode[codepos:]));
			break;
		}
		formpos = formpos+1;
		codepos = codepos+1;
	}
	var ret []string = eval_expr(expr,env);
	return ret
}

//should evaluate the code
func eval_expr(code []string, env *Environment) ([]string){
	if len(code) == 0{
		return []string{};
	}

	var oper string = code[0];

	if val,ok := env.check(oper); ok && len(code) == 1{
		return val;
	} else if len(code) == 1{
		return code;
	}

	if ret, ok := env.apply_forms(code); ok{
		return ret;
	}

	if oper == "eval" {
		return eval_expr(code[1:],env);
	} else if oper == "evalis" {
		var ret []string;
		for _,arg := range deparenth(code[1:]){
			ret = eval_expr(arg,env);
		}
		return ret;
	} else if oper == "eval_env" || oper == "void"{
		args := deparenth(code[1:]);
		tobind_vars := deparenth(args[0]);
		void_env := newEnvironment();
		for _,tobind := range tobind_vars{
			tmp := deparenth(tobind);
			void_env.update(tmp[0][0],eval_expr(tmp[1],env));
		}
		return eval_expr(args[1],void_env);
	} else if oper == "evalx"{
		return eval_expr(eval_expr(code[1:],env),env)
	} else if oper == "escape"{
		var toadd map[string][]string;
		args := deparenth(code[1:]);
		tobind_vars := deparenth(args[0]);
		for _,tobind := range tobind_vars{
			tmp := deparenth(tobind);
			toadd[tmp[0][0]] = eval_expr(tmp[1],env);
		}

		hold := env.environments[0];
		env.contract();

		for key,val := range toadd{
			env.update(key,val);
		}

		ret := eval_expr(args[1],env);

		env.extend();
		for key,val := range hold{
			env.update(key,val);
		}
		return ret;
	} else if oper == "print"{
			var toprint string = "";
			var params [][]string = deparenth(code[1:]);
			for _,param := range params[:len(params)-1]{
				toprint = toprint+listCombinef(eval_expr(param,env))+" ";
			}
			fmt.Print(toprint+listCombinef(eval_expr(params[len(params)-1],env)));
			return []string{""};
	} else if oper == "println"{
			var toprint string = "";
			var params [][]string = deparenth(code[1:]);
			for _,param := range params[:len(params)-1]{
				toprint = toprint+listCombinef(eval_expr(param,env))+" ";
			}
			fmt.Println(toprint+listCombinef(eval_expr(params[len(params)-1],env)));
			return []string{""};
	} else if oper == "printall"{
			var params [][]string = deparenth(code[1:]);
			for _,param := range params{
				fmt.Println(listCombinef(eval_expr(param,env)));
			}
			return []string{""};
	} else if oper == "global"{
		var toadd map[string][]string;
		args := deparenth(code[1:]);
		tobind_vars := deparenth(args[0]);
		for _,tobind := range tobind_vars{
			tmp := deparenth(tobind);
			toadd[tmp[0][0]] = eval_expr(tmp[1],env);
		}

		var hold []map[string][]string;
		for len(env.environments) > 1{
			tmp := env.environments[0];
			hold = append([]map[string][]string{},tmp);
			env.contract();
		}

		for key,val := range toadd{
			env.update(key,val);
		}

		ret := eval_expr(args[1],env);

		for _,hold_env := range hold{
			env.extend();
			for key,val := range hold_env{
				env.update(key,val);
			}
		}
		return ret;
	} else if oper == "q" || oper == "quote"{
		return code[1:];
	} else if oper == "parenth"{
		return append([]string{"("},append(code[1:],")")...);
	} else if oper == "rel" || oper == "release"{
		var ret []string = make([]string,0);
		for _,tok := range code[1:]{
			if _,ok := env.check(tok); ok{
				ret = append(ret,env.get(tok)...);
			} else{
				ret = append(ret,tok);
			}
		}
		return ret;
	} else if oper == "free"{
		return free(code[1:],env);
	} else if oper == "rel-lam"{
		args1 := deparenth(code[1:])
		return release_lambda(eval_expr(args1[0],env),args1[1:],env);
	} else if oper == "body"{
		test := listSplit(eval_expr(deparenth(code[1:])[0],env),"|");
		if len(test) != 3{
			return test[0];
		} else{
			return test[2];
		}
	} else if oper == "lambda"{
		tmp := deparenth(code[1:]);
		return append([]string{"procedure","|"},append(tmp[0],append([]string{"|"},tmp[1]...)...)...);
	} else if oper == "fexpr"{
		tmp := deparenth(code[1:]);
		return append([]string{"fprocedure","|"},append(tmp[0],append([]string{"|"},tmp[1]...)...)...);
	} else if oper == "def" || oper == "define"{
		var varname string = code[1];
		env.update(varname,eval_expr(deparenth(code[2:])[0],env));
		return []string{};
	} else if oper == "form"{
		args := deparenth(code[1:]);
		env.set_forms(listCombinef(args[0]),args[1]);
		return []string{};
	} else if oper == "cond"{
		for _,clause := range deparenth(code[1:]){
			var case_statement [][]string = deparenth(clause)
			if eval_expr(case_statement[0],env)[0] == "#t"{
				ret := eval_expr(case_statement[1],env);
				return ret;
			}
		}
		return []string{};
	} else if oper == "loop"{
		var args [][]string = deparenth(code[1:]);
		iter,err := strconv.Atoi(eval_expr(args[0],env)[0]);
		handle(err);
		i := 0;
		var ret []string;
		for i < iter{
			for _,arg := range args[1:]{
				ret = eval_expr(arg,env);
			}
			i = i+1;
		}
		return ret;
	} else if oper == "pass"{
		ret := make([]string,0);
		args := deparenth(code[1:]);
		for _,arg := range args{
			ret = append(ret,eval_expr(arg,env)...);
		}
		return ret;
	} else if oper == "len" || oper == "length"{
		return []string{strconv.Itoa(len(code)-1)};
	} else if oper == "exit"{
			exit_code,err := strconv.Atoi(code[1]);
			handle(err);
			os.Exit(exit_code);
			return []string{};
	}else{
		code2 := deparenth(code);
		var proc []string = eval_expr(code2[0],env);
		var args [][]string = make([][]string,0);
		if listSplit(proc,"|")[0][0] == "fprocedure"{
			args = code2[1:];
		} else{
			for  _,val := range code2[1:]{
				toadd := eval_expr(val,env);
				args = append(args,toadd);
			}
		}
		var check string = listSplit(proc,"|")[0][0];
		if((check != "native") && (check != "procedure") && (check != "fprocedure")){
			err_string := listCombine(proc)+" is not a procedure";
			//panic(err_string);
			return []string{err_string};
		}
		return apply(proc,args,env);
	}
}
