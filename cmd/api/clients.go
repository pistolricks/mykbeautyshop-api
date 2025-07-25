package main

import (
	"net/http"
)

func (app *application) listClientsHandler(w http.ResponseWriter, r *http.Request) {

	clients, metadata, err := app.models.Clients.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"clients": clients, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
