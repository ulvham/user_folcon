package main

import (
	"database/sql"
	//"errors"
	"fmt"
	"io/ioutil"
	//"net"
	"net/http"
	"os"
	"path/filepath"
	_ "pkg/github.com/mattn/go-sqlite3"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
	//"net/url"
)

type Page struct {
	NameFile string
	Title    string
	Body     string
	Menu     string
	Menuw    string
	Menue    string
	IsCheck  string
	HidMenu  string
	CurUser  string
	ExitP    string
	BaseName string
	AChecks  string
	ASelects string
	Fnames   string
}

var curuser = make([]string, 0)

func cur_dir() string {
	cwd, _ := os.Getwd() //Получить текущую директорию
	//fmt.Println(cwd)
	return cwd

}

func (p *Page) save(r *http.Request) error {

	var _, ttt = findcookie(r, "username")

	_ = os.Mkdir("Pages", 0755)
	_ = os.Mkdir("Pages/"+ttt, 0755)
	_ = os.Mkdir("Pages/"+ttt+"/Title", 0755)
	_ = os.Mkdir("Pages/"+ttt+"/Text", 0755)
	_ = os.Mkdir("Pages/"+ttt+"/BaseN", 0755)
	filename1 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Title/" + p.NameFile + ".title"
	filename2 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Text/" + p.NameFile + ".text"
	filename3 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/BaseN/" + ttt + ".bn"
	visibl := "@_no"
	if p.IsCheck == "on" {
		visibl = "@yes"
	}
	ioutil.WriteFile(filename1, []byte(p.Title+visibl), 0600) // права на запись и на чтение
	ioutil.WriteFile(filename3, []byte(p.BaseName), 0600)     // права на запись и на чтение
	return ioutil.WriteFile(filename2, []byte(p.Body), 0600)  // права на запись и на чтение
}

func loadbasename(namefile string, r *http.Request) string {
	var _, ttt = findcookie(r, "username")
	filename3 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/BaseN/" + ttt + ".bn"
	bname, _ := ioutil.ReadFile(filename3)
	pathdb := strings.Split(string(bname), "|")
	table_begin := "<select onclick = \"gogo()\" size=\"1\" id=\"tags\" name=\"tags\"><option>Выберите базу</option>"
	table_end := "</select>"
	table_col_begin := "<option value="
	table_col_end := "</option>"
	sss := table_begin

	for i := 1; i <= len(pathdb)-2; i++ {
		sss += table_col_begin
		sss += "\"" + pathdb[i] + "\">" + pathdb[i]
		sss += table_col_end

	}
	sss += table_end
	return sss
}

func delfiles(w http.ResponseWriter, r *http.Request) {

	var _, ttt = findcookie(r, "username")

	dirname := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Filez/"
	d, err := os.Open(dirname)
	if err != nil {
		return
	}
	fi, err := d.Readdir(-1)
	if err != nil {
		return
	}
	fz := r.FormValue("delete")
	fzz := strings.Split(fz, "|")

	for _, fi := range fi {
		for i := 0; i < len(fzz); i++ {

			if fi.Name() == fzz[i] {
				fmt.Println(fi.Name(), fzz[i], "/Pages/"+ttt+"/Filez/"+fi.Name())
				_ = os.Remove("." + string(filepath.Separator) + "/Pages/" + ttt + "/Filez/" + fi.Name())
			}
		}

	}

	fnames := loadfilesname(r)

	fmt.Fprintf(w, fnames)

}

func loadfilesname(r *http.Request) string {
	var _, ttt = findcookie(r, "username")
	var fileall []string
	dirname := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Filez/"
	d, err := os.Open(dirname)
	if err != nil {
		return "not found"
	}
	fi, err := d.Readdir(-1)
	if err != nil {
		return "not found"
	}

	for _, fi := range fi {

		//filename := fi.Name()

		fileall = append(fileall, "<p><input class=\"ch1\" name=\""+fi.Name()+"\" type=\"checkbox\" ><a href=\"../../files_"+ttt+"/"+fi.Name()+"\" id=\"button\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-radio-off\"></span>"+fi.Name()+"</a></p>")
	}

	sort.Strings(fileall)
	fileall = append(fileall, "<p><input name=\"allcheck1\" onclick=\"chall('allcheck1','ch1');\" type=\"checkbox\" ><a href=\"#\" id=\"button\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-radio-off\"></span>Выбрать/Убрать ВСЕ</a></p><div><center><input type=\"submit\" onclick = \"delfiles()\" value=\"Удалить отмеченные\"></center></div>")

	fileall_ := strings.Join(fileall, " ")
	return fileall_

}

func loadPage(namefile string, r *http.Request) (*Page, error) {

	var _, ttt = findcookie(r, "username")

	filename1 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Title/" + namefile + ".title"
	filename2 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Text/" + namefile + ".text"
	//filename3 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/BaseN/" + namefile + ".bn"
	title, err1 := ioutil.ReadFile(filename1)
	body, err2 := ioutil.ReadFile(filename2)
	//bname, err3 := ioutil.ReadFile(filename3)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	// if err3 != nil {
	// 	return nil, err3
	// }
	var title1 string
	title1 = strings.Replace(strings.Replace(string(title), "@_no", "", -1), "@yes", "", -1)

	var menu, menu0 []string
	var menuw, menue string

	dirname := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Title/"
	d, err := os.Open(dirname)
	if err != nil {
	}
	fi, err := d.Readdir(-1)
	if err != nil {
	}
	len0 := len(title)
	hidm := ""
	for _, fi := range fi {
		lenf0 := len(".title")
		lenf1 := len(fi.Name())
		lenf2 := lenf1 - lenf0

		filename := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Title/" + fi.Name()
		titlemenu, _ := ioutil.ReadFile(filename)
		lent := len(titlemenu)
		title2 := strings.Replace(strings.Replace(string(titlemenu), "@_no", "", -1), "@yes", "", -1)

		if string(titlemenu)[lent-4:] == "@yes" {
			if fi.Name()[:lenf2] == "index" {
				if fi.Name()[:lenf2] == namefile {
					menu0 = append(menu0, "<p><a href=\"../view/"+fi.Name()[:lenf2]+"\""+"id=\"button\" class=\"ui-state-default ui-corner-all ui-state-hover\"><span class=\"ui-icon ui-icon-bullet\"></span>"+string(title2)+"</a></p>")
				}

				if fi.Name()[:lenf2] != namefile {
					menu0 = append(menu0, "<p><a href=\"../view/"+fi.Name()[:lenf2]+"\""+"id=\"button\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-radio-off\"></span>"+string(title2)+"</a></p>")
				}
			}

			if fi.Name()[:lenf2] != "index" {
				if fi.Name()[:lenf2] == namefile {
					menu = append(menu, "<p><a href=\"../view/"+fi.Name()[:lenf2]+"\""+"id=\"button\" class=\"ui-state-default ui-corner-all ui-state-hover\"><span class=\"ui-icon ui-icon-bullet\"></span>"+string(title2)+"</a></p>")
				}

				if fi.Name()[:lenf2] != namefile {
					menu = append(menu, "<p><a href=\"../view/"+fi.Name()[:lenf2]+"\""+"id=\"button\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-radio-off\"></span>"+string(title2)+"</a></p>")
				}
			}

		} else {

			if fi.Name()[:lenf2] == namefile {
				//hidm = hidm + "<script type=\"text/javascript\"> jQuery(document).ready(function(){$(\"#accordion\").accordion(\"option\", \"active\", 2);});</script><p><a href=\"../view/" + fi.Name()[:lenf2] + "\"" + "id=\"button\" class=\"ui-state-default ui-corner-all ui-state-hover\"><span class=\"ui-icon ui-icon-bullet\"></span>" + string(title2) + "</a></p>"
				hidm = hidm + "<p><a href=\"../view/" + fi.Name()[:lenf2] + "\"" + "id=\"button\" class=\"ui-state-default ui-corner-all ui-state-hover\"><span class=\"ui-icon ui-icon-bullet\"></span>" + string(title2) + "</a></p>"
			}

			if fi.Name()[:lenf2] != namefile {
				hidm = hidm + "<p><a href=\"../view/" + fi.Name()[:lenf2] + "\"" + "id=\"button\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-radio-off\"></span>" + string(title2) + "</a></p>"
			}
		}
	}

	sort.Strings(menu)
	menu2 := append(menu0, menu...)

	menu1 := strings.Join(menu2, " ")
	menuw = "\"../viewpage/" + namefile + "\""
	menue = "\"../edit/" + namefile + "\""

	isch := ""
	if string(title)[len0-4:] == "@yes" {
		isch = "checked"
	}
	extp := ""

	////_, ttt := findcookie(r, "username")
	if ttt != "" {
		extp = "<h3><a href=\"#\">Файлы</a></h3><div><p><a href=\"../../files_" + ttt + "\" id=\"button\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-power\"></span>Filez (" + ttt + ")</a></p></div><h3><a href=\"#\">Выйти</a></h3><div><p><a href=\"#\" id=\"dialog_link1\" class=\"ui-state-default ui-corner-all\"><span class=\"ui-icon ui-icon-power\"></span>Выйти (" + ttt + ")</a></p><div id=\"dialog1\" title=\"Покинуть\"><p>Вы действительно хотите выйти?</p></div></div>"
	}
	//var checks, selects string
	bname := loadbasename(namefile, r)
	fnames := loadfilesname(r)

	return &Page{NameFile: namefile, Title: title1, Body: string(body), Menu: menu1, Menuw: menuw, Menue: menue, IsCheck: isch, HidMenu: hidm, CurUser: ttt, ExitP: extp, BaseName: bname, Fnames: fnames}, nil
}

//const lenPath = len("/view/")

var templates = make(map[string]*template.Template)

var namefileValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func viewHandler(w http.ResponseWriter, r *http.Request, namefile string) {
	if r.FormValue("new") != "" {
		newpath := r.FormValue("new")
		http.Redirect(w, r, "/edit/"+newpath, http.StatusFound)
		return
	}
	p, err := loadPage(namefile, r)
	if err != nil {
		http.Redirect(w, r, "/edit/"+namefile, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func viewHandler1(w http.ResponseWriter, r *http.Request, namefile string) {
	p, err := loadPage(namefile, r)
	if err != nil {
		http.Redirect(w, r, "/edit/"+namefile, http.StatusFound)
		return
	}
	renderTemplate(w, "viewpage", p)
}

func deltHandler(w http.ResponseWriter, r *http.Request, namefile string) {

	var _, ttt = findcookie(r, "username")

	_ = os.Remove("Pages/" + ttt + "/Title/" + namefile + ".title")
	_ = os.Remove("Pages/" + ttt + "/Text/" + namefile + ".text")
	//_ = os.Remove("Pages/" + ttt + "/BaseN/" + namefile + ".bn")
	http.Redirect(w, r, "/view/index", http.StatusFound)
}

func viewHandlerDef(w http.ResponseWriter, r *http.Request, namefile string) {
	p, err := loadPage(namefile, r)
	if err != nil {
		http.Redirect(w, r, "/edit/index", http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, namefile string) {
	p, err := loadPage(namefile, r)
	if err != nil {
		p = &Page{NameFile: namefile}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, namefile string) {
	body := r.FormValue("body")
	title := r.FormValue("title")
	isch := r.FormValue("visib")
	bname := r.FormValue("b_n")
	p := &Page{NameFile: namefile, Title: title, Body: body, IsCheck: isch, BaseName: bname}
	err := p.save(r)
	if err != nil {
		http.Error(w, "Error!", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+namefile, http.StatusFound)
}

func init() {
	for _, tmpl := range []string{"edit", "view", "viewpage"} {
		templates[tmpl], _ = template.ParseFiles("Makets/" + tmpl + ".html")
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates[tmpl].Execute(w, p) //
	if err != nil {
		http.Error(w, "Error!", http.StatusInternalServerError)
	}
}

//func getnamefile(w http.ResponseWriter, r *http.Request) (namefile string, err error) {
//	namefile = r.URL.Path[lenPath:]
//	if !namefileValidator.MatchString(namefile) {
//		http.NotFound(w, r)
//		err = errors.New("Invalid Page namefile")
//	}
//	return
//}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namefile := ""
		var _, ttt = findcookie(r, "username")

		if ttt != "" {
			curuser = adduser(curuser, ttt)
		}

		var lenpath int
		//strings.Split(r.URL.Path, "/")[1] != "unlogon"
		if strings.Split(r.URL.Path, "/")[1] != "upload" && strings.Split(r.URL.Path, "/")[1] != "unlogon" && strings.Split(r.URL.Path, "/")[1] != "ot4" {
			if !strings.Contains(r.URL.Path, "/favicon.ico") {
				lenpath = len(strings.Split(r.URL.Path, "/")[1]) + 2
			} else {
				lenpath = 1
			}

		} else {
			lenpath = 1
		}

		if strings.Contains(r.URL.Path, "/favicon.ico") {
			//tttt := "/files_" + ttt + "/"
			if r.URL.Path == "/files_"+ttt {
				//fmt.Println(r.URL.Path)
				http.Redirect(w, r, "/files_"+ttt+"/", http.StatusFound)
				//http.NotFound(w, r)
				//return
			}
			if !strings.Contains(r.URL.Path, "/files_"+ttt) {
				fmt.Println(r.URL.Path)
				http.Redirect(w, r, "/files_"+ttt+"/index.html", http.StatusFound)
				//http.NotFound(w, r)
				return
			}
		}

		if r.URL.Path != "/" {
			if finduser(curuser, ttt) {
				//fmt.Println(len(strings.Split(r.URL.Path, "/")[1]))
				//namefile = r.URL.Path[len(strings.Split(r.URL.Path, "/")[1])+2:]
				fmt.Println(lenpath, "", r.URL.Path)
				namefile = r.URL.Path[lenpath:]

			} else {
				namefile = "login"
				if autoriz(r.FormValue("login"), r.FormValue("password")) == true {
					r.SetBasicAuth(r.FormValue("login"), r.FormValue("password"))
					cookie := http.Cookie{Name: "username", Value: r.FormValue("login"), Path: "/"}
					http.SetCookie(w, &cookie)
					curuser = adduser(curuser, r.FormValue("login"))
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"), curuser, "enter:", r.FormValue("login"))
					//namefile = r.URL.Path[len(strings.Split(r.URL.Path, "/")[1])+2:]
					namefile = r.URL.Path[lenpath:]
					//					err := os.Mkdir(r.FormValue("login"), 0755)
					//					if err == nil {
					//						http.Handle("/files"+r.FormValue("login")+"/", http.StripPrefix("/files"+r.FormValue("login")+"/", http.FileServer(http.Dir(r.FormValue("login")))))
					//					}
					http.Redirect(w, r, "/view/index", http.StatusFound)
				}
			}
		}

		if r.URL.Path == "/" {
			if finduser(curuser, ttt) {
				namefile = "index"
			} else {
				namefile = "login"
				if autoriz(r.FormValue("login"), r.FormValue("password")) == true {
					r.SetBasicAuth(r.FormValue("login"), r.FormValue("password"))
					namefile = "index"
					cookie := http.Cookie{Name: "username", Value: r.FormValue("login"), Path: "/"}
					http.SetCookie(w, &cookie)
					curuser = adduser(curuser, r.FormValue("login"))
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"), curuser, "enter:", r.FormValue("login"))
					//					err := os.Mkdir(r.FormValue("login"), 0755)
					//					if err == nil {
					//						http.Handle("/files"+r.FormValue("login")+"/", http.StripPrefix("/files"+r.FormValue("login")+"/", http.FileServer(http.Dir(r.FormValue("login")))))
					//					}
					http.Redirect(w, r, "/view/index", http.StatusFound)
				}
			}
		}
		if !namefileValidator.MatchString(namefile) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, namefile)
	}
}

func adduser(slice []string, data string) []string {
	if !finduser(slice, data) {
		l := len(slice)
		newSlice := make([]string, (l + 1))
		for i, c := range slice {
			newSlice[i] = c
		}

		slice = newSlice[0 : l+1]
		slice[l] = data
	}
	return slice
}

func deluser(slice []string, data string) []string {
	if finduser(slice, data) {
		l := len(slice)
		newSlice := make([]string, (l))
		var ii = 0
		for _, c := range slice {
			if c != data {
				newSlice[ii] = c
				ii++
			}
		}
		slice = newSlice[0 : l-1]
	}
	return slice
}

func finduser(slice []string, data string) bool {
	for _, c := range slice {
		if c == data {
			return true
		}
	}
	return false
}

func findcookie(r *http.Request, zn string) (int, string) {
	for i, cookie := range r.Cookies() {
		if cookie.Name == zn {
			return i, cookie.Value
		}
	}
	return -1, ""
}

func autoriz(login string, password string) bool {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		//fmt.Println(err)
		return false
	}
	defer db.Close()

	//просто запросы, возможно делать включения сюда

	//	rows, err := db.Query("select id, name, email from users where name = ?")
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	defer rows.Close()
	//	for rows.Next() {
	//		var id int
	//		var name string
	//		var email string
	//		rows.Scan(&id, &name, &email)
	//		println(id, name, email)
	//	}
	//	rows.Close()

	// многофункциональное препарирование

	stmt, err := db.Prepare("select password from users where name = ?")
	if err != nil {
		//fmt.Println(err)
		return false
	}
	defer stmt.Close()
	var pwd string
	err = stmt.QueryRow(login).Scan(&pwd)
	if err != nil {
		//fmt.Println(err)
		return false
	}
	if password != pwd {
		return false
	}
	return true

}

func entry() {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer db.Close()

	//просто запросы, возможно делать включения сюда

	rows, err := db.Query("select name from users")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		os.Mkdir("./Pages/"+name+"/Filez/", 0755)
		//if err == nil {
		go http.Handle("/files_"+name+"/", http.StripPrefix("/files_"+name+"/", http.FileServer(http.Dir("./Pages/"+name+"/Filez/"))))
		//}
		println(name)
	}
	rows.Close()
}

func uni(database string, tablename string, r *http.Request) ([][]string, string) {

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, "nil"
	}
	defer db.Close()
	countrows, err := db.Query("select count(*) from " + tablename)
	if err != nil {
		fmt.Println(err)
		return nil, "nil"
	}

	var ii int //---------------------------
	for countrows.Next() {
		countrows.Scan(&ii)
		//fmt.Println(ii)
	}
	countrows.Close()

	columnsdef, err := db.Query("select * from " + tablename)
	if err != nil {
		fmt.Println(err)
		return nil, "nil"
	}
	columns, _ := columnsdef.Columns()
	idn := 1
	lll := len(columns)
	vv := make([]string, lll)
	//-------------------$("#content").hasClass("newsBlock")
	colzn := "<tr id=\"t0\">"
	for izn, zn := range columns {
		if r.FormValue(zn) == "on" {
			vv[izn] = zn
			colzn += "<td id=\"t" + strconv.Itoa(idn) + "\" onclick=\"if ($('.s" + strconv.Itoa(idn) + "').hasClass('fixx')) {$('.s" + strconv.Itoa(idn) + "').removeClass('fixx');} else {$('.s" + strconv.Itoa(idn) + "').addClass('fixx');};\">" + zn + "</td>"
			idn++

		} else {
			vv[izn] = "@@@nil@@@"
		}

	}
	colzn += "</tr>"
	//-------------------
	fmt.Println(vv)
	l := lll
	//l := len(columns)

	newSlice := make([][]string, l)

	iii := 0
	inew := 0
	var ic int
	for i := range newSlice {
		if vv[i] != "@@@nil@@@" {
			if inew == 0 {
				ic = 0
			} else {
				ic = i
			}
			newSlice[ic] = make([]string, ii)
			rows, err := db.Query("select " + columns[i] + " from " + tablename)

			fmt.Println(columns[i])
			fmt.Println(ic, i, inew)
			if err != nil {
				fmt.Println(err)
				return nil, "nil"
			}
			defer rows.Close()
			var row string
			for rows.Next() {
				rows.Scan(&row)
				newSlice[ic][iii] = row
				//fmt.Println(row)
				iii++

			}
			rows.Close()
			iii = 0
			inew = 1
		}

	}
	return newSlice, colzn

}

func univtable(vt [][]string, ht string) *string {
	var text string
	// table_begin := "<table class=\"table_uni\" border=2px cellspacing=0>"
	// table_end := "</table>"

	// table_row_begin := "<tr>"
	// table_row_end := "</tr>"

	// table_col_begin := "<td>"
	// table_col_end := "</td>"
	text += "<table class=\"table_uni\" border=2px cellspacing=0>"
	text += ht
	idnn := 1
	idnnnn := 0
	for iii := 0; iii < len(vt[0]); iii++ {
		if idnn <= len(vt) && idnnnn == 0 {
			text += "<tr id=\"s0\" onclick=\"if ($(this).hasClass('fixx')) {$(this).removeClass('fixx');} else {$(this).addClass('fixx');};\">"
		} else {
			text += "<tr onclick=\"if ($(this).hasClass('fixx')) {$(this).removeClass('fixx');} else {$(this).addClass('fixx');};\">"
		}
		for i := range vt {
			for ii := range vt[i] {
				if idnn <= len(vt) && idnnnn == 0 {
					text += "<td class=\"s" + strconv.Itoa(idnn) + "\" id=\"s" + strconv.Itoa(idnn) + "\">"
				} else {
					text += "<td class=\"s" + strconv.Itoa(idnn) + "\">"
				}
				idnn++

				text += vt[i][ii+iii]
				text += "</td>"
				break
			}

		}
		text += "</tr>"
		idnn = 1
		idnnnn++
	}
	text += "</table>"
	return &text
}

func univinput(vt [][]string) *string {

	var text string
	table_begin := "<label for=\"tablename\"><b>Выберите имя таблицы: </b></label><select onclick = \"go()\" size=\"1\" id=\"tablename\" name=\"tablename\"><option disabled>Выберите таблицу</option>"
	table_end := "</select>"
	table_col_begin := "<option value="
	table_col_end := "</option>"
	text += table_begin
	for iii := 0; iii < len(vt[0]); iii++ {
		for i := range vt {
			for ii := range vt[i] {
				text += table_col_begin
				text += "\"" + vt[i][ii+iii] + "\">" + vt[i][ii+iii]
				text += table_col_end
				break
			}
		}

	}
	text += table_end
	return &text
}

func uni0(database string) [][]string {

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil
	}
	defer db.Close()
	countrows, err := db.Query("SELECT count(*) FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var ii int //---------------------------
	for countrows.Next() {
		countrows.Scan(&ii)
		//fmt.Println(ii)
	}
	countrows.Close()

	l := 1

	newSlice := make([][]string, l)

	iii := 0
	for i := range newSlice {
		newSlice[i] = make([]string, ii)
		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer rows.Close()
		var row string
		for rows.Next() {
			rows.Scan(&row)
			newSlice[i][iii] = row
			//fmt.Println(row)
			iii++

		}
		rows.Close()
		iii = 0
	}
	return newSlice

}

func universal_o(w http.ResponseWriter, r *http.Request, namefile string) {
	//r.FormValue("fullpathname")

	// pathdb := strings.Split(r.FormValue("fullpathname"), "|")[0]
	// tabledb := strings.Split(r.FormValue("fullpathname"), "|")[1]

	pathdb := r.FormValue("tags") + ".db"
	tabledb := r.FormValue("tablename")

	query_lowpassa := *univtable(uni(pathdb, tabledb, r))

	var _, ttt = findcookie(r, "username")

	_ = os.Mkdir("Pages", 0755)
	_ = os.Mkdir("Pages/"+ttt, 0755)
	_ = os.Mkdir("Pages/"+ttt+"/Title", 0755)
	_ = os.Mkdir("Pages/"+ttt+"/Text", 0755)
	filename1 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Title/" + "otchet" + time.Now().Format("20060102150405") + ".title"
	filename2 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/Text/" + "otchet" + time.Now().Format("20060102150405") + ".text"
	ioutil.WriteFile(filename1, []byte("Отчёт "+time.Now().Format("02-01-2006(15:04)")+"@yes"), 0600) // права на запись и на чтение
	ioutil.WriteFile(filename2, []byte(query_lowpassa), 0600)                                         // права на запись и на чтение
	// многофункциональное препарирование

	path__ := strings.Split(filename1, "/")[len(strings.Split(filename1, "/"))-1]
	path_ := strings.Split(path__, ".")[0]
	//fmt.Println(path_)

	filename3 := "." + string(filepath.Separator) + "/Pages/" + ttt + "/BaseN/" + path_ + ".bn"
	ioutil.WriteFile(filename3, []byte(time.Now().Format("02-01-2006(15:04)")), 0600)
	//path_:=strings.Split(path__,"@")[1]
	//fmt.Println(path_)
	http.Redirect(w, r, "/view/"+path_, http.StatusFound)

}

func checkboxtable(database string, tablename string) *string {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil
	}
	defer db.Close()

	columnsdef, err := db.Query("select * from " + tablename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	columns, _ := columnsdef.Columns()
	//-------------------$("#content").hasClass("newsBlock")
	colzn := "<br><br><label><b>Выберите столбцы: </b></label><br><table><tr></td><td><td><input name=\"allcheck0\" onclick=\"chall('allcheck0','ch0');\" type=\"checkbox\" checked></td></tr>"
	for _, zn := range columns {
		colzn += "<tr><td>" + zn + "</td>" + "<td><input class=\"ch0\" name=\"" + zn + "\" type=\"checkbox\" checked></td></tr>"
	}
	colzn += "</table>"
	return &colzn
}

func UploadServer(w http.ResponseWriter, req *http.Request, namefile string) {
	var _, ttt = findcookie(req, "username")
	//fz := req.FormValue("path")
	//fmt.Println(fz)
	if ttt != "" {
		file, handler, err := req.FormFile("userfile")
		//fmt.Println(handler.Filename)
		if err != nil {
			fmt.Println(err)
			fnames := loadfilesname(req)
			fmt.Fprintf(w, fnames)
			//http.Redirect(w, req, "/#tabs-5", http.StatusFound)
			return
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			fnames := loadfilesname(req)
			fmt.Fprintf(w, fnames)
			//http.Redirect(w, req, "/#tabs-5", http.StatusFound)
			return
		}
		filename := "./Pages/" + ttt + "/Filez/" + handler.Filename
		err = ioutil.WriteFile(filename, data, 0777)
		if err != nil {
			fmt.Println(err)
			fnames := loadfilesname(req)
			fmt.Fprintf(w, fnames)
			//http.Redirect(w, req, "/#tabs-5", http.StatusFound)
			return
		}
		fnames := loadfilesname(req)
		fmt.Fprintf(w, fnames)
		//http.Redirect(w, req, "/#tabs-5", http.StatusFound)
	} else {
		fnames := loadfilesname(req)
		fmt.Fprintf(w, fnames)

		//http.Redirect(w, req, "/#tabs-5", http.StatusFound)
	}
}

func byselect(w http.ResponseWriter, req *http.Request) {
	//tablename := req.FormValue("selection")
	dbname := req.FormValue("databaze")
	fmt.Println(dbname)
	selects := fmt.Sprintf(*univinput(uni0(dbname + ".db")))
	//selects_ :=
	//fmt.Sprintf(*checkboxtable(dbname+".db",tablename))
	//checks = *checkboxtable(pathdb, tabledb)
	//resstring :=fmt.Sprintf("<p>The id is %s</p>","moocow");
	fmt.Fprintf(w, selects)

	//io.WriteString(http.Conn, resstring);
}

func bytab(w http.ResponseWriter, req *http.Request) {
	tablename := req.FormValue("selection")
	dbname := req.FormValue("databaze")
	selects := fmt.Sprintf(*checkboxtable(dbname+".db", tablename))
	fmt.Fprintf(w, selects)
	// tablename := req.FormValue("selection")
	//selects := fmt.Sprintf(*univinput(uni0(pathdb)),*checkboxtable("predpr1.db",tablename))

	//checks = *checkboxtable(pathdb, tabledb)
	//resstring :=fmt.Sprintf("<p>The id is %s</p>","moocow");
	// fmt.Fprintf(w, selects)

	//io.WriteString(http.Conn, resstring);
}

func unlogon(w http.ResponseWriter, r *http.Request, namefile string) {

	_, ttt := findcookie(r, "username")
	curuser = deluser(curuser, ttt)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), curuser, "exit:", ttt)
	cookie := http.Cookie{Name: "username", Value: "", Path: "/"}
	http.SetCookie(w, &cookie)

	//http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "files")
	//})

	////dsm := *http.DefaultServeMux
	//map[string]string{"key1": "val1", "key2": "val2"}
	//fmt.Println(dsm)

	//dsm[t http.ServeMux] = _, false
	////fmt.Println(dsm)

	//http.Server{
	//	Handler:     http.FileServer(http.Dir("files")),
	//	ReadTimeout: 30 * time.Second,
	//}

	http.Redirect(w, r, "/view/index", http.StatusFound)
}

func main() {

	go http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	go http.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("pic"))))

	entry()

	//	for i := 0; i < 10000; i++ {
	//		ii := strconv.Itoa(i)
	//		go http.Handle("/files"+ii+"/", http.StripPrefix("/files"+ii+"/", http.FileServer(http.Dir("name1"))))
	//		fmt.Println(ii)
	//	}
	//go http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("name1"))))
	//go http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("name2"))))

	go http.HandleFunc("/view/", makeHandler(viewHandler))
	go http.HandleFunc("/viewpage/", makeHandler(viewHandler1))
	go http.HandleFunc("/edit/", makeHandler(editHandler))
	go http.HandleFunc("/save/", makeHandler(saveHandler))
	go http.HandleFunc("/delete/", makeHandler(deltHandler))
	go http.HandleFunc("/unlogon", makeHandler(unlogon))
	go http.HandleFunc("/ot4", makeHandler(universal_o))
	go http.HandleFunc("/", makeHandler(viewHandlerDef))
	go http.HandleFunc("/upload", makeHandler(UploadServer))
	go http.HandleFunc("/byselect", byselect)
	go http.HandleFunc("/bytab", bytab)
	go http.HandleFunc("/delfiles", delfiles)

	//srv := &http.Server{
	//						ReadTimeout: 30 * time.Second,
	//					 }

	go http.ListenAndServe("0.0.0.0:8080", nil)

	//ch := make(chan bool)
	//<-ch

	select {}

}
