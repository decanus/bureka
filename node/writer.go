package node

type Writer struct {

}

//func (w *Writer) Send(ctx context.Context, target state.Peer, msg pb.Message) error {
//	out, ok := w.writers[peer.ID(target)]
//	if !ok {
//		return fmt.Errorf("peer %s not found", string(target))
//	}
//
//	out <- msg
//	return nil
//}
//
