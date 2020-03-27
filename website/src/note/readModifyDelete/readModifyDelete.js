import React from "react";
import '../note.css'
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import Spinner from "react-bootstrap/Spinner";
import ReadModifyDeleteForm from "./readModifyDeleteForm";

class Show extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showReadPasswordModal: true,
            showModifyDeleteModal: false,
            modifyDeleteModal: {
                title: '',
                button: ''
            },
            readDeleteError: false,
            modifyError: false,
            passwordForm: React.createRef(),
            contentRef: React.createRef(),
            ttlRef: React.createRef(),
            readDeleteSpinner: false,
            modifySpinner: false,
            note: {
                content: '',
                ttl: ''
            },
            chars: 140,
            writeAccess: false,
        }
    }

    onReadPasswordSubmit = (event) => {
        event.preventDefault();
        this.getNote(this.props.match.params.id, this.state.passwordForm.current.value)
    };

    passwordFocus = () => {
        this.setState({readDeleteError: false, modifyError: false});
    };

    contentChange = (event) => {
        event.persist();
        this.setState(() => ({
            chars: 140 - event.target.value.length
        }));
    };

    onWriteAccessCheckChange = () => {
        this.setState((state) => ({writeAccess: !state.writeAccess}))
    };

    closeModifyDeleteModal = () => {
        if (this.state.modifyDeleteModal.button === 'Close') {
            this.setState(() => ({showModifyDeleteModal: false}))
        } else {
            this.props.history.push("/notes")
        }
    };

    onModify = (event) => {
        event.preventDefault();
        let data = {
            content: this.state.contentRef.current.value,
            ttl: parseInt(this.state.ttlRef.current.value)
        };
        this.modifyNote(data, this.props.match.params.id, this.state.passwordForm.current.value)
    };

    onDelete = (event) => {
        event.preventDefault();
        this.deleteNote(this.props.match.params.id, this.state.passwordForm.current.value)
    };

    async getNote(id, password) {
        this.setState({readDeleteSpinner: true});
        let response;
        try {
            response = await this.get(id, password)
        } catch (err) {
            console.error(err);
            this.setState({readDeleteError: true, readDeleteSpinner: false});
            return
        }
        this.setState({
                note: {content: response.content, ttl: response.ttl},
                chars: 140 - response.content.length,
                readDeleteSpinner: false,
                readDeleteError: false,
                showReadPasswordModal: false
            },
        );
    }

    async modifyNote(data, id, password) {
        this.setState({modifySpinner: true});
        try {
            await this.put(data, id, password)
        } catch (err) {
            this.setState({modifyError: true, modifySpinner: false});
            return
        }
        this.setState({
            modifyDeleteModal: {
                title: 'Note modified!',
                button: 'Close'
            },
            showModifyDeleteModal: true,
            modifySpinner: false
        });
    }

    async deleteNote(id, password) {
        this.setState({readDeleteSpinner: true});
        try {
            await this.delete(id, password)
        } catch (err) {
            this.setState({modifyError: true, readDeleteSpinner: false});
            return
        }
        this.setState({
            modifyDeleteModal: {
                title: 'Note deleted!',
                button: 'Create another'
            },
            showModifyDeleteModal: true,
            readDeleteSpinner: false
        });
    }

    async get(id, password) {
        let response = await fetch(
            '/api/notes/' + id,
            {
                headers: {
                    'Password': password,
                }
            }
        );
        if (response.status !== 200) {
            throw new Error(response.statusText)
        }
        return await response.json()
    }

    async put(data, id, password) {
        let response = await fetch(
            '/api/notes/' + id,
            {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Password': password
                },
                body: JSON.stringify(data)
            }
        );
        if (response.status !== 200) {
            throw new Error(response.statusText)
        }
    }

    async delete(id, password) {
        let response = await fetch(
            '/api/notes/' + id,
            {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'Password': password
                },
            }
        );
        if (response.status !== 200) {
            throw new Error(response.statusText)
        }
    }

    render() {
        return (
            <>
                <ReadModifyDeleteForm
                    defaultContent={this.state.note.content}
                    defaultTTL={this.state.note.ttl}
                    contentRef={this.state.contentRef}
                    contentChange={this.contentChange}
                    ttlRef={this.state.ttlRef}
                    chars={this.state.chars}
                    checkChange={this.onWriteAccessCheckChange}
                    writeAccessChecked={this.state.writeAccess}
                    deleteSpinner={this.state.readDeleteSpinner}
                    modifySpinner={this.state.modifySpinner}
                    onDelete={this.onDelete}
                    onModify={this.onModify}
                    onPasswordFocus={this.passwordFocus}
                    modifyError={this.state.modifyError}
                    passwordRef={this.state.passwordForm}
                />
                <Modal show={this.state.showReadPasswordModal} dialogClassName="passwordModal">
                    <Modal.Header closeButton>
                        <Modal.Title>Type password</Modal.Title>
                    </Modal.Header>
                    <Form onSubmit={this.onReadPasswordSubmit}>
                        <Modal.Body>
                            <Form.Group controlId="passwordForm">
                                <Form.Control type="password" ref={this.state.passwordForm} maxLength={50}
                                              onFocus={this.passwordFocus} placeholder="read or write password"/>
                            </Form.Group>
                            {this.state.readDeleteError && <Form.Text className="passwordHint">
                                incorrect password or note does not exist
                            </Form.Text>}
                        </Modal.Body>
                        <Modal.Footer>
                            <Button type="submit"
                                    variant='primary'
                                    disabled={this.state.readDeleteSpinner}>
                                {this.state.readDeleteSpinner && <Spinner
                                    as="span"
                                    animation="grow"
                                    size="sm"
                                    role="status"
                                    aria-hidden="true"
                                />}
                                Send
                            </Button>
                        </Modal.Footer>
                    </Form>
                </Modal>
                <Modal show={this.state.showModifyDeleteModal}>
                    <Modal.Header closeButton>
                        <Modal.Title>{this.state.modifyDeleteModal.title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Footer>
                        <Button variant="primary" onClick={this.closeModifyDeleteModal}>
                            {this.state.modifyDeleteModal.button}
                        </Button>
                    </Modal.Footer>
                </Modal>
            </>
        )
    }
}

export default Show