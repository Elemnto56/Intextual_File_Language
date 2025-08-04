package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TUI() {
	app := tview.NewApplication()
	// Info page
	info_box := tview.NewTextView().SetText(`
                                              %%                                             
                                           %%%%%                                             						
                                        %%%%%%%%                                             
                                     %%%%%%%%%%%%%%%                                         				     	--- Intext TUI Version --------------------------			
                                    %%%%%%%%%%%%%%%%%%%%%%%%@                                                       Version       : 1.0     
                                       %%%%%%%%%%    %%%%%%%%%%%%                            						Codename      : Simplicity 
                         %%               %%%%%%            %%%%%%%%%                        						Release Date  : TBD
                      %%%%%%                 %%%                 %%%%%%                      						Interface Type: Terminal UI 
                   %%%%%%                                           %%%%%%                   						Built With    : tview
                 %%%%%%                                               %%%%%%                 						Status        : Early, yet Stable
               %%%%%                                                     %%%%%               
              %%%%                                                         %%%%%             
            %%%%%                %%%%                                       %%%%%            
           %%%%                %%%%%%%%%%%%%%%%%%%%%%%                        %%%%           						--- Intext Core Version ------------------------
          %%%%               %%%% @%%%           %%%%%%%                       %%%%%         						Version       : v0.8		
        %%%%@              %%%%     %%%        %%%%   %%%                        %%%%        						Codename      : Where's the Logic?	
       %%%%              %%%%         %%%    %%%%       %%%                       %%%%       						Release Date  : TBD
       %%%             %%%             %%%%%%%%          %%%%                      %%%%      						Runtime Engine: ISEC
      %%%%           %%%                 %%%%              %%%                     %%%%      				
     %%%%          %%%                                      %%%%                    %%%%     			
     %%%        %%%%                              %%%         %%%                    %%%     
    %%%%      %%%%                 %%           %%%%            %%%                  %%%%    						--- Documentation --------------------------------
    %%%     %%%%                 %%%%         %%%%               %%%%                 %%%   						Link          : TODO
   %%%%      %%%               %%%          %%%%                   %%%                %%%%   						Last Updated  : TBD			
   %%%%        %%%          @%%%          %%%%          %%%         %%%%              %%%%   						CLI Help      : Run with "--help" or "-H" for options
   %%%          %%%%        %%          %%%%          %%%@            %%%              %%%   						
   %%%            %%%                 %%%%          %%%                 %%%            %%%   						
   %%%             %%%%             %%%%          %%%          %%        %%%%                
   %%%               %%%          %%%%         %%%%          %%%%          %%%               
   %%%%                %%%                   %%%%          %%%%             %%%%             
   %%%%                 %%%@               %%%%          %%%%                %%%%            
    %%%                   %%%            %%%%          %%%%                %%%%              
    %%%%                   %%%%         %%%          %%%%                %%%%                
     %%%                     %%%                   %%%@                %%%%           %%%    
     %%%%                      %%%               %%%                 %%%%          %%%%%%    
      %%%%     %                %%%%           %%%                 %%%%         %%%%%%%%%    
       %%%% %%%%                  %%%         %%                 %%%%       %%%%%%%%%%%%%    
       %%%%%%%%%                   %%%%                        %%%@          %%%%%%%%%%%%    
      %%%%%%%%%%                     %%%                     %%%                %%%%%%%%%    
   %%%%%%%%%%%%%                       %%%                 %%%                  %%%%%%%%%    
      %%%%%%%%%%                        %%%%             %%%                  %%%%%   %%%    
         %%%%%%%                          %%%         %%%%                   %%%%%           
            %%%%                           %%%%     %%%%                   %%%%%             
               %                             %%%  %%%%                   %%%%%               
                                               %%%%%                   %%%%%%                
                                                %%                   %%%%%%                  
                        %%%                                       %%%%%%                     
                       %%%%%%%%                               %%%%%%%                        
                           %%%%%%%%%                     %%%%%%%%%                           
                               %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%                               
                                     %%%%%%%%%%%%%%%%%%%                                     
	`)

	info_box.SetBorder(true).SetBorderAttributes(tcell.AttrBlink).SetBorderColor(tcell.ColorBlue).SetTitle("Intext Information Page").SetTitleAlign(tview.AlignLeft)

	// Run Page
	input := tview.NewInputField().SetFieldWidth(30).SetLabel("Input Filename")

	runner := tview.NewForm().
		AddFormItem(input).
		AddButton("Submit", func() {
		})

	runner.SetBorder(true).SetTitle("Run Menu").SetBorderColor(tcell.ColorGreen)

	// Set up
	list := tview.NewList().
		AddItem("Run", "Run an Intext File", 'r', nil).
		AddItem("Info", "See version info", 'i', nil).
		AddItem("Exit", "Leave the TUI", 'q', func() { app.Stop() })

	pages := tview.NewPages().
		AddPage("menu", list, true, true).
		AddPage("Run", runner, true, false).
		AddPage("Info", info_box, true, false)

	pages.SetTitle("Intext's TUI")
	// Functionatily
	list.SetSelectedFunc(func(index int, main string, secondary string, shortcut rune) {
		switch main {
		case "Run":
			pages.SwitchToPage("Run")
		case "Info":
			pages.SwitchToPage("Info")
		case "Exit":
			app.Stop()
		}
	})

	info_box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage("menu")
		}
		return nil
	})

	runner.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage("menu")
		}
		return nil
	})

	app.SetRoot(pages, true)
	err := app.Run()
	if err != nil {
		panic("Something bad happened")
	}
}
