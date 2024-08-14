
interface SSEResponse {
    done:  boolean;
    source:  String;
    data: String | Object;
    duration: String;
    success: boolean;
}
