import React from "react";
import Form from "react-bootstrap/Form";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import Button from "react-bootstrap/Button";
import Spinner from "react-bootstrap/Spinner";

const NoteForm = (props) => {
    return (
        <div className="characters-hint">
            <Form onSubmit={props.onSubmit} onChange={props.onChange}>
                <Form.Group controlId="content">
                    <Form.Control ref={props.content} as="textarea" maxLength={140} minLength={1} required rows={6}
                                  placeholder="note content"/>
                    <Form.Text className="text-muted" float="right">{props.chars} characters left</Form.Text>
                </Form.Group>
                <OverlayTrigger placement="right" delay={{show: 250, hide: 400}} overlay={
                    <Tooltip id="button-tooltip">if empty your note will NOT be secured!</Tooltip>}>
                    <Form.Group controlId="readPassword">
                        <Form.Control type="password" maxLength={50} placeholder="read password"/>
                    </Form.Group>
                </OverlayTrigger>
                <OverlayTrigger placement="right" delay={{show: 250, hide: 400}} overlay={
                    <Tooltip id="button-tooltip">allows modifying content and TTL after note creation</Tooltip>}>
                    <Form.Group controlId="writePassword">
                        <Form.Control type="password" maxLength={50} placeholder="write password"/>
                    </Form.Group>
                </OverlayTrigger>
                <OverlayTrigger placement="right" delay={{show: 250, hide: 400}} overlay={
                    <Tooltip id="button-tooltip">default & maximum 1 day</Tooltip>}>
                    <Form.Group controlId="ttl">
                        <Form.Control ref={props.ttl} type="number" min={1} max={1440} placeholder="minutes to expire"/>
                    </Form.Group>
                </OverlayTrigger>
                <Button variant="primary" type="submit" disabled={props.spinner}>
                    {props.spinner && <Spinner
                        as="span"
                        animation="grow"
                        size="sm"
                        role="status"
                        aria-hidden="true"/>}
                    Create</Button>
            </Form>
        </div>)
};

export default NoteForm
