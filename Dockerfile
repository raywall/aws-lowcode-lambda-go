FROM public.ecr.aws/lambda/provided:al2

# Copia apenas o necessário
COPY bin/bootstrap ${LAMBDA_TASK_ROOT}

CMD ["bootstrap"]