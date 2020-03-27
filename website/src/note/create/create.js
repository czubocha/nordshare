import React from 'react';
import NoteForm from "./createForm";
import Modal from 'react-bootstrap/Modal'
import Button from 'react-bootstrap/Button'

class Create extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            chars: 140,
            content: '',
            readPassword: '',
            writePassword: '',
            ttl: 1440,
            showModal: false,
            showSpinner: false,
            modal: {
                title: '',
                content: ''
            },
            copyButton: {
                text: 'Copy',
                variant: 'primary'
            },
            noteLink: React.createRef()
        };
    }

    handleChange = (event) => {
        if (event.target.id === 'content') {
            this.setState({chars: 140 - event.target.value.length})
        }
        this.setState({[event.target.id]: event.target.value});
    };

    handleSubmit = (event) => {
        event.preventDefault();
        let data = {
            content: this.state.content,
            readPassword: this.state.readPassword,
            writePassword: this.state.writePassword,
            ttl: parseInt(this.state.ttl)
        };
        this.sendNote(data)
    };

    refreshState = () => {
        this.setState({
            showModal: false,
            showSpinner: false,
            modal: {
                title: '',
                content: ''
            },
            copyButton: {
                text: 'Copy',
                variant: 'primary'
            },
            noteLink: React.createRef()
        })
    };

    toggleSpinner = () => {
        this.setState((state) => ({showSpinner: !state.showSpinner}))
    };

    async post(data) {
        let response = await fetch(
            '/api/notes',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            }
        );
        if (response.status !== 201) {
            throw new Error(response.statusText)
        }
        return await response.json()
    }

    async sendNote(data) {
        this.toggleSpinner();
        let response;
        try {
            response = await this.post(data)
        } catch (err) {
            this.setState({
                modal: {title: 'Error', content: 'Something went wrong...'},
                showModal: true,
                showSpinner: false,
            });
            return
        }
        this.setState({
            modal: {
                title: 'Note created!',
                content: window.location.protocol + '//' + window.location.host + '/notes/' + response.id,
            },
            showModal: true,
            showSpinner: false,
        });
    }

    async copyCodeToClipboard() {
        if (navigator.clipboard === undefined) {
            console.error('clipboard is undefined');
            this.setState({copyButton: {text: "Can't copy, sorry...", variant: 'danger'}});
            return
        }
        try {
            await navigator.clipboard.writeText(this.state.noteLink.current.getAttribute('href'));
        } catch (err) {
            console.error(err);
            this.setState({copyButton: {text: "Can't copy, sorry...", variant: 'danger'}})
        }
        this.setState({copyButton: {text: 'Copied!', variant: 'success'}})
    };

    render() {
        return (
            <><NoteForm chars={this.state.chars} onChange={this.handleChange}
                        onSubmit={this.handleSubmit} spinner={this.state.showSpinner}/>
                <br/>
                <Modal show={this.state.showModal}>
                    <Modal.Header closeButton><Modal.Title>{this.state.modal.title}</Modal.Title></Modal.Header>
                    <Modal.Body>
                        <a ref={this.state.noteLink}
                           href={this.state.modal.content}>{this.state.modal.content}</a></Modal.Body>
                    <Modal.Footer>
                        <Button variant="secondary" onClick={this.refreshState}>Create one more</Button>
                        {!(this.state.modal.title === 'Error') &&
                        <Button onClick={() => this.copyCodeToClipboard()} variant={this.state.copyButton.variant}>
                            {this.state.copyButton.text}</Button>}
                        {this.state.copyButton.variant === 'danger' &&
                        <p className="text-muted security">security reasons.</p>}
                    </Modal.Footer>
                </Modal></>)
    }
}

export default Create