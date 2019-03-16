// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

Cypress.Commands.add('actionStart', () => {
  cy.server();
  cy.route({
    method: 'POST',
    url: '/*',
  }).as('formAction');
});

Cypress.Commands.add('actionWait', () => {
  cy.wait('@formAction');
});

Cypress.Commands.add('resetCMS', () => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:8080/',
    form: true,
    body: {
      action: 'index-cypress-reset',
    }
  }).then((response) => {
    expect(response.body).to.eq('OK');
  });
});

Cypress.Commands.add('installCMS', () => {
  cy.actionStart();
  cy.resetCMS();

  cy.visit('http://localhost:8080/cp/');
  cy.get('.form-signin input[name=name]').type('fave');
  cy.get('.form-signin input[name=user]').type('root');
  cy.get('.form-signin input[name=password]').type('root');
  cy.get('.form-signin button').click();
  cy.actionWait();

  cy.visit('http://localhost:8080/cp/');
  cy.get('.form-signin input[name=first_name]').type('First');
  cy.get('.form-signin input[name=last_name]').type('Last');
  cy.get('.form-signin input[name=email]').type('example@example.com');
  cy.get('.form-signin input[name=password]').type('example@example.com');
  cy.get('.form-signin button').click();
  cy.actionWait();
});

Cypress.Commands.add('loginCMS', () => {
  cy.actionStart();
  cy.visit('http://localhost:8080/cp/');
  cy.get('.form-signin input[name=email]').type('example@example.com');
  cy.get('.form-signin input[name=password]').type('example@example.com');
  cy.get('.form-signin button').click();
  cy.actionWait();
});
