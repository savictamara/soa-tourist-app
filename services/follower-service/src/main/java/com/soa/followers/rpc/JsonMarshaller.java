package com.soa.followers.rpc;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.grpc.MethodDescriptor;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.UncheckedIOException;

public class JsonMarshaller<T> implements MethodDescriptor.Marshaller<T> {
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();
    private final Class<T> type;

    public JsonMarshaller(Class<T> type) {
        this.type = type;
    }

    @Override
    public InputStream stream(T value) {
        try {
            return new ByteArrayInputStream(OBJECT_MAPPER.writeValueAsBytes(value));
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }

    @Override
    public T parse(InputStream stream) {
        try {
            return OBJECT_MAPPER.readValue(stream, type);
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }
}
